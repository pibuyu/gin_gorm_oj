package main

import (
	"gin_gorm_o/router"
)

func main() {
	r := router.Router()
	r.Run(":8080")
}
