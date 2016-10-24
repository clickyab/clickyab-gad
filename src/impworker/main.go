package main

import (
	"config"
	"models"
	"rabbit"
	"transport"
	"utils"
	"version"

	"redis"

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	version.PrintVersion().Info("Application started")
	models.Initialize()
	rabbit.Initialize()
	aredis.Initialize()
	//
	//tmp := transport.Impression{
	//	User: <-utils.ID,
	//	ImpID: <-utils.ID,
	//	AdID: 12,
	//	CampaignID: 3456789,
	//	UserAgent: <-utils.ID,
	//	WinnerBID: 1,
	//	Status: 1,
	//	Cookie:false,
	//	Suspicious:false,
	//
	//}
	//
	//s, _ := json.MarshalIndent(tmp, "\t", "\t")
	//fmt.Println(string(s))

	exit := make(chan chan struct{})

	go func() {
		err := rabbit.RunWorker(
			config.Config.AMQP.Exchange,
			"cy.imp",
			"cy_imp_queue",
			&transport.Impression{},
			impWorker,
			10,
			exit,
		)
		if err != nil {
			// Fatal is only allowed in main
			logrus.Fatal(err)
		}
	}()

	utils.WaitSignal(exit)
	logrus.Info("goodbye")
}
