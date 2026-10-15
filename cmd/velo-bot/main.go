package main

import (
	"fmt"

	"github.com/Samswrld02/Velo-Discord-Bot/bot"
	"github.com/Samswrld02/Velo-Discord-Bot/logger"
	"github.com/Samswrld02/Velo-Discord-Bot/types"
	"github.com/joho/godotenv"
)



func main() {
	fmt.Printf("Velo-Helpdesk is running: ") 

	//load in env
	godotenv.Load()

	//setup config struct
	config := types.NewConfig()

	logger, file, err := logger.NewLogger()
	if err != nil {
		panic("Error creating logger: " + err.Error())
	}
	defer file.Close()

	bot := bot.NewBot(logger, config)
	bot.Run()
}

