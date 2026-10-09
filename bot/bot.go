package bot

import "github.com/bwmarrin/discordgo"


var BotToken string


func run() {

	//create a discord session
	discord, err := discordgo.New("Bot " + BotToken)

	//event handler
	discord.AddHandler()
}

func handleConnection() {
	
}