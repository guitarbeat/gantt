package main

import "honnef.co/go/tools/config"

type Page struct {
	Cfg     Config
	Modules Modules
}

type Modules []Module
type Module struct {
	Cfg  config.Config
	Tpl  string
	Body interface{}
}
