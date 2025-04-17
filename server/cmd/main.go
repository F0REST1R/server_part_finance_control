package main

import (
	"Server_part_finance_control/server/pkg/config"
	"Server_part_finance_control/server/pkg/logger"
	"fmt"
)

func main() {
	var err error
	cfg, err := config.New()
	if err != nil{
		fmt.Println("Error")
	}

	log := logger.SetupLogger(cfg.ENV)

	log.Info("starting application")

	// if err != nil{
	// 	log.Error("error message: ", slog.String("error", err.Error()))
	// }
}