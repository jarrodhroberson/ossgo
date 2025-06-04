package valkey

import (
	"fmt"
)

var DefaultDatabase = Database{Name: "default", Id: 0}

type Database struct {
	Name string
	Id   int
}

func (d Database) String() string {
	return fmt.Sprintf("%d:%s", d.Id, d.Name)
}

type KeyFunc func(string) string