package main

import (
	"GOIM/router"
)

func main() {

	r := router.Router()
	r.Run(":8080")
}
