package entity

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Menu struct {
	Uuid uuid.UUID `db:"uuid" json:"uuid"`
	Name string    `db:"name" json:"name"`
	Url  string    `db:"url" json:"url"`
}

type MenuWithItems struct {
	Uuid  uuid.UUID  `db:"uuid" json:"uuid"`
	Name  string     `db:"name" json:"name"`
	Url   string     `db:"url" json:"url"`
	Items []MenuItem `json:"items"`
}

type MenuItem struct {
	Uuid      uuid.UUID `db:"uuid" json:"uuid"`
	ShortName string    `db:"short_name" json:"short_name"`
	Name      string    `db:"name" json:"name"`
	Price     int       `db:"price" json:"price"`
	MenuUuid  uuid.UUID `db:"menu_uuid" json:"menu_uuid"`
}

type NewMenu struct {
	Name  string        `json:"name"`
	Url   string        `json:"url"`
	Items []NewMenuItem `json:"items"`
}

type NewMenuItem struct {
	ShortName string `json:"short_name"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
}

func (*Repository) GetAllMenus(tx *sqlx.Tx) ([]MenuWithItems, error) {
	menusMap := map[uuid.UUID]*MenuWithItems{}
	rows, err := tx.Queryx("SELECT * FROM menus")
	if err != nil {
		return nil, fmt.Errorf("could not get all menus from db: %w", err)
	}
	for rows.Next() {
		var menu MenuWithItems
		rows.StructScan(&menu)
		menusMap[menu.Uuid] = &menu
	}

	rows, err = tx.Queryx("SELECT mi.* FROM menus m JOIN menu_items mi on m.uuid = mi.menu_uuid")
	if err != nil {
		return nil, fmt.Errorf("could not get all menu_items from db: %w", err)
	}
	for rows.Next() {
		var menuItem MenuItem
		rows.StructScan(&menuItem)
		menusMap[menuItem.MenuUuid].Items = append(menusMap[menuItem.MenuUuid].Items, menuItem)
	}

	menus := make([]MenuWithItems, 0, len(menusMap))
	for _, value := range menusMap {
		menus = append(menus, *value)
	}

	return menus, nil
}

func (repo *Repository) GetMenu(tx *sqlx.Tx, menuUuid uuid.UUID) (*MenuWithItems, error) {
	var menu Menu
	if err := tx.Get(&menu, "SELECT * FROM menus WHERE uuid = $1", menuUuid); err != nil {
		return nil, fmt.Errorf("failed to get menu %s: %w", menuUuid, err)
	}

	var items []MenuItem
	if err := tx.Select(&items, "SELECT * FROM menu_items WHERE menu_uuid = $1", menuUuid); err != nil {
		return nil, fmt.Errorf("failed to get menu items for menu %s: %w", menuUuid, err)
	}

	return &MenuWithItems{
		Uuid:  menu.Uuid,
		Name:  menu.Name,
		Url:   menu.Url,
		Items: items,
	}, nil
}

func (repo *Repository) GetMenuItem(tx *sqlx.Tx, menuItemUuid uuid.UUID) (*MenuItem, error) {
	var menuItem MenuItem
	if err := tx.Get(&menuItem, "SELECT * FROM menu_items WHERE uuid = $1", menuItemUuid); err != nil {
		return nil, fmt.Errorf("failed to get menu item %s: %w", menuItemUuid, err)
	}

	return &menuItem, nil
}

func (repo *Repository) CreateMenu(tx *sqlx.Tx, menu *NewMenu) (*MenuWithItems, error) {
	var uuidString string
	err := tx.Get(&uuidString, "INSERT INTO menus (name, url) VALUES ($1, $2) RETURNING uuid", menu.Name, menu.Url)
	if err != nil {
		return nil, fmt.Errorf("could not create menu %s: %w", menu.Name, err)
	}
	createdMenuUuid := uuid.Must(uuid.FromString(uuidString))

	menuItems := []MenuItem{}
	for _, menuItem := range menu.Items {
		var createdMenuItem MenuItem
		err = tx.Get(
			&createdMenuItem,
			"INSERT INTO menu_items (name, short_name, price, menu_uuid) VALUES ($1, $2, $3, $4) RETURNING uuid, name, short_name, price, menu_uuid",
			menuItem.Name, menuItem.ShortName, menuItem.Price, createdMenuUuid,
		)
		if err != nil {
			return nil, fmt.Errorf("could not create menuItem %s: %w", menuItem.Name, err)
		}
		menuItems = append(menuItems, createdMenuItem)
	}

	return &MenuWithItems{Uuid: createdMenuUuid, Name: menu.Name, Url: menu.Url, Items: menuItems}, nil
}

func (repo *Repository) UpdateMenu(tx *sqlx.Tx, menuUuid uuid.UUID, menu *NewMenu) (*MenuWithItems, error) {
	existingMenu, err := repo.GetMenu(tx, menuUuid)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE menus SET name = $2, url = $3 WHERE uuid = $1", menuUuid, menu.Name, menu.Url)
	if err != nil {
		return nil, fmt.Errorf("could not update menu %s: %w", menuUuid, err)
	}

	existingMenuItemMap := make(map[string]MenuItem)
	for _, item := range existingMenu.Items {
		existingMenuItemMap[item.ShortName] = item
	}

	menuItemMap := make(map[string]NewMenuItem)
	for _, newItem := range menu.Items {
		menuItemMap[newItem.ShortName] = newItem
		if _, found := existingMenuItemMap[newItem.ShortName]; !found {
			_, err = repo.CreateMenuItem(tx, &newItem, menuUuid)
			if err != nil {
				return nil, fmt.Errorf("error while creating missing menu items: %w", err)
			}
		}
	}

	for _, existingItem := range existingMenu.Items {
		if _, found := menuItemMap[existingItem.ShortName]; !found {
			if err = repo.DeleteMenuItem(tx, existingItem.Uuid); err != nil {
				return nil, fmt.Errorf("error while deleting orphan menu items: %w", err)
			}
		}
	}

	return repo.GetMenu(tx, menuUuid)
}

func (repo *Repository) CreateMenuItem(tx *sqlx.Tx, menuItem *NewMenuItem, menuUuid uuid.UUID) (*MenuItem, error) {
	var createdMenuItem MenuItem
	err := tx.Get(
		&createdMenuItem,
		"INSERT INTO menu_items (name, short_name, price, menu_uuid) VALUES ($1, $2, $3, $4) RETURNING uuid, name, short_name, price, menu_uuid",
		menuItem.Name, menuItem.ShortName, menuItem.Price, menuUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create menu item %s: %w", menuItem.ShortName, err)
	}

	return &createdMenuItem, nil
}

func (repo *Repository) DeleteMenuItem(tx *sqlx.Tx, menuItemUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM menu_items WHERE uuid = $1", menuItemUuid)
	if err != nil {
		return fmt.Errorf("could not delete menu item %s: %w", menuItemUuid, err)
	}

	return nil
}

func (repo *Repository) DeleteMenu(tx *sqlx.Tx, menuUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM menus WHERE uuid = $1", menuUuid)
	if err != nil {
		return fmt.Errorf("could not delete menu %s: %w", menuUuid, err)
	}

	return nil
}
