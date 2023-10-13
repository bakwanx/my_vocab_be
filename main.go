package main

import (
	"my_vocab/config"
	"my_vocab/route"
)

func main() {
	config.InitDb()
	route.InitRoute()
}
