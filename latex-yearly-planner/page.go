package main

type PageModule struct {
	Cfg     Config
	Modules Modules
}

type Modules []Module
type Module struct {
	Cfg  Config
	Tpl  string
	Body interface{}
}
