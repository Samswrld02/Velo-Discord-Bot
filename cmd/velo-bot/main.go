package main

import (
	"fmt"

	"github.com/Samswrld02/Velo-Discord-Bot/bot"
	"github.com/joho/godotenv"
)



func main() {
	fmt.Printf("Velo-Helpdesk is running: ") 

	//load in env
	godotenv.Load()

	bot := bot.NewBot()
	bot.Run()
}

