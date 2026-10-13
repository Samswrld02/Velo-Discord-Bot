package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/Samswrld02/Velo-Discord-Bot/types"

	redisconnection "github.com/Samswrld02/Velo-Discord-Bot/redis-connection"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}


func (b *Bot) Run() {

	//open session
	b.session.Open()
	defer b.session.Close()

	
	//event handler
	// discord.AddHandler(newMessage)
	// discord.AddHandler(readEventFromQueue)
	b.ticketInQueue()


	fmt.Println("bot is running")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}


func handleTicketEvent() {

}

func NewBot() *Bot{
	BotToken, found := os.LookupEnv("DISCORD_BOT_TOKEN")

	if found != true {
		panic("Error with the bot token")
	}

	//create a discord session
	discord, err := discordgo.New("Bot " +  BotToken)

	if err != nil {
		panic(err)
	}
	
	return &Bot{session: discord}
}

//method for checking if ticket event was put into the queue
func (b *Bot) ticketInQueue() (bool, error) {

	//get env
	queue, found := os.LookupEnv("REDIS_QUEUE")
	if found != true {
		panic("Error with given queue")
	}

	//make redis connection
	r := redisconnection.NewRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PW"))

	for  {
		//blocking function -> synchronous for now
		value := r.ReadFromQueue(queue)
		b.handleData(value)
	}
	
}

//method for processing data into struct
func (b *Bot) handleData(data []byte) {
	//import ticket struct
	var ticket types.Ticket;

	//unmarshal data into struct
	json.Unmarshal(data, &ticket)

	ch, err := b.session.UserChannelCreate(ticket.Discord_id)

	if err != nil {
		fmt.Println("Error occurred: " + err.Error())
	}

	_, err = b.session.ChannelMessageSendComplex(ch.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed {
			{
				Title: "Ticket assigned by " + ticket.Assigned_by,
				Description: ticket.Subject,
			},

		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "view ticket",
						Style: discordgo.LinkButton,
						URL: ticket.Url,
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}


