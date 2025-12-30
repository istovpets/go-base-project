package rest

import "github.com/go-fuego/fuego"

func (r *Rest) addHandlers() {
	fuego.Get(r.srv, "/ping", pingHandler)
}
