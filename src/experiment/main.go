package main

import (
	"config"
	"rabbit"
	"time"
	"transport"
)

func main() {
	config.Initialize()

	rabbit.Initialize()
	defer rabbit.FinalizeWait()

	rabbit.PublishAfter("cy.click", transport.Click{}, time.Second)
}
