package main

import (
	"ticket-office/internal/config"
)

func main() {
	cfg := config.MustLoad()

	_ = cfg
}
