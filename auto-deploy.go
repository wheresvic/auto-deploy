package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/wheresvic/auto-deploy/src/adconfiguration"
	"github.com/wheresvic/auto-deploy/src/adserver"
	"github.com/wheresvic/auto-deploy/src/adversion"
)

func main() {
	log.Print("Starting auto-deploy")
	rand.Seed(time.Now().UnixNano())

	adVersion := adversion.GetCurrentVersion()
	log.Print(adVersion)
	initConfig, err := adconfiguration.LoadAndSetConfiguration("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	server := adserver.InitServer(initConfig, adVersion)
	adserver.Start(server, initConfig.Server.HTTPPort)
}
