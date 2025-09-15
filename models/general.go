// Package models all structs used in the project
package models

import (
	"fmt"
	"html/template"
	"sync"
)

type PageData struct {
	Theme string
	Title string
	Data  *struct{}
}

type User struct {
	Username, FirstName, LastName, Email, Hash, Token string
}

type Message struct {
	Message string
}

type Templates struct {
	m sync.Map
}

func (t *Templates) Load(name string) (*template.Template, error) {
	val, ok := t.m.Load(name)
	if !ok {
		return nil, fmt.Errorf("was not able to parse template: %s", name)
	}

	if tmpl, ok := val.(*template.Template); ok {
		return tmpl, nil
	} else {
		return nil, fmt.Errorf("was not able to typecast: %s", name)
	}
}

func (t *Templates) Add(key string, template *template.Template) {
	t.m.Store(key, template)
}

func (t *Templates) List() {
	t.m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}
