package model

import "github.com/gofrs/uuid"

type Menu struct {
	UUID *uuid.UUID
	Name string
	URL  string
}

type MenuWithItems struct {
	UUID  *uuid.UUID
	Name  string
	URL   string
	Items []MenuItem
}

type MenuItem struct {
	UUID      *uuid.UUID
	ShortName string
	Name      string
	Price     int
	MenuUUID  *uuid.UUID
}
