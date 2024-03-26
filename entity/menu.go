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
	menus_map := map[uuid.UUID]*MenuWithItems{}
	rows, err := tx.Queryx("SELECT * FROM menus")
	if err != nil {
		return nil, fmt.Errorf("could not get all menus from db: %w", err)
	}
	for rows.Next() {
		var menu MenuWithItems
		rows.StructScan(&menu)
		menus_map[menu.Uuid] = &menu
	}

	rows, err = tx.Queryx("SELECT mi.* FROM menus m JOIN menu_items mi on m.uuid = mi.menu_uuid")
	if err != nil {
		return nil, fmt.Errorf("could not get all menu_items from db: %w", err)
	}
	for rows.Next() {
		var menu_item MenuItem
		rows.StructScan(&menu_item)
		menus_map[menu_item.MenuUuid].Items = append(menus_map[menu_item.MenuUuid].Items, menu_item)
	}

	menus := make([]MenuWithItems, 0, len(menus_map))
	for _, value := range menus_map {
		menus = append(menus, *value)
	}

	return menus, nil
}

func (repo *Repository) GetMenu(tx *sqlx.Tx, menu_uuid uuid.UUID) (*MenuWithItems, error) {
	var menu Menu
	if err := tx.Get(&menu, "SELECT * FROM menus WHERE uuid = $1", menu_uuid); err != nil {
		return nil, fmt.Errorf("failed to get menu %s: %w", menu_uuid, err)
	}

	var items []MenuItem
	if err := tx.Select(&items, "SELECT * FROM menu_items WHERE menu_uuid = $1", menu_uuid); err != nil {
		return nil, fmt.Errorf("failed to get menu items for menu %s: %w", menu_uuid, err)
	}

	return &MenuWithItems{
		Uuid:  menu.Uuid,
		Name:  menu.Name,
		Url:   menu.Url,
		Items: items,
	}, nil
}

func (repo *Repository) GetMenuItem(tx *sqlx.Tx, menu_item_uuid uuid.UUID) (*MenuItem, error) {
	var menuItem MenuItem
	if err := tx.Get(&menuItem, "SELECT * FROM menu_items WHERE uuid = $1", menu_item_uuid); err != nil {
		return nil, fmt.Errorf("failed to get menu item %s: %w", menu_item_uuid, err)
	}

	return &menuItem, nil
}

func (repo *Repository) CreateMenu(tx *sqlx.Tx, menu *NewMenu) (*MenuWithItems, error) {
	var uuid_string string
	err := tx.Get(&uuid_string, "INSERT INTO menus (name, url) VALUES ($1, $2) RETURNING uuid", menu.Name, menu.Url)
	if err != nil {
		return nil, fmt.Errorf("could not create menu %s: %w", menu.Name, err)
	}
	created_menu_uuid := uuid.Must(uuid.FromString(uuid_string))

	menu_items := []MenuItem{}
	for _, menu_item := range menu.Items {
		var created_menu_item MenuItem
		err = tx.Get(
			&created_menu_item,
			"INSERT INTO menu_items (name, short_name, price, menu_uuid) VALUES ($1, $2, $3, $4) RETURNING uuid, name, short_name, price, menu_uuid",
			menu_item.Name, menu_item.ShortName, menu_item.Price, created_menu_uuid,
		)
		if err != nil {
			return nil, fmt.Errorf("could not create menu_item %s: %w", menu_item.Name, err)
		}
		menu_items = append(menu_items, created_menu_item)
	}

	return &MenuWithItems{Uuid: created_menu_uuid, Name: menu.Name, Url: menu.Url, Items: menu_items}, nil
}

func (repo *Repository) UpdateMenu(tx *sqlx.Tx, menu_uuid uuid.UUID, menu *NewMenu) (*MenuWithItems, error) {
	existing_menu, err := repo.GetMenu(tx, menu_uuid)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE menus SET name = $2, url = $3 WHERE uuid = $1", menu_uuid, menu.Name, menu.Url)
	if err != nil {
		return nil, fmt.Errorf("could not update menu %s: %w", menu_uuid, err)
	}

	existingMenuItemMap := make(map[string]MenuItem)
	for _, item := range existing_menu.Items {
		existingMenuItemMap[item.ShortName] = item
	}

	menuItemMap := make(map[string]NewMenuItem)
	for _, newItem := range menu.Items {
		menuItemMap[newItem.ShortName] = newItem
		if _, found := existingMenuItemMap[newItem.ShortName]; !found {
			_, err = repo.CreateMenuItem(tx, &newItem, menu_uuid)
			if err != nil {
				return nil, fmt.Errorf("error while creating missing menu items: %w", err)
			}
		}
	}

	for _, existingItem := range existing_menu.Items {
		if _, found := menuItemMap[existingItem.ShortName]; !found {
			if err = repo.DeleteMenuItem(tx, existingItem.Uuid); err != nil {
				return nil, fmt.Errorf("error while deleting orphan menu items: %w", err)
			}
		}
	}

	return repo.GetMenu(tx, menu_uuid)
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
