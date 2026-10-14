package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"

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
	b.handleQueue()


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
func (b *Bot) handleQueue() {

	//retrieve redis db for the queue
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("error with db .env variable")
	}
	
	//create redis connection
	r := redisconnection.NewRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PW"), db)

	//channel for the queued events
	eventChannel := make(chan []byte)

	//worker that reads data from the redis queue into the channel
	b.readEventWorker(eventChannel, r)

	//specified ammount of workers that process data from the channel concurrently
	b.handleTicketWorkers(5, eventChannel)
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

	message := createMessage(ticket)

	_, err = b.session.ChannelMessageSendComplex(ch.ID, message)

	if err != nil {
		fmt.Println(err)
	}
}

//function for creating the message component/content
func createMessage(ticket types.Ticket) *discordgo.MessageSend {
	return &discordgo.MessageSend{
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
	}
}

//workers that handle incoming data into eventchannel
func (b *Bot) handleTicketWorkers(workers int, eventChannel <-chan []byte ) {
	//worker pool that handles the jobs from the queue
	for i := 0; i < workers; i++ {
		//dispatch worker
		go func(){
			//read from channel, blocks if channel is empty
			for job := range eventChannel {
				b.handleData(job)
			}
		}()
	}
}

//processo/worker of redis queue
func (b *Bot) readEventWorker(eventChannel chan []byte, r *redisconnection.RedisC) {
		queue, found := os.LookupEnv("REDIS_QUEUE")
		if found != true {
			panic("Error with given queue")
		}
		//get env
		go func() {		
		//loop infinitely blpop will block if queue is empty
			for {
				value := r.ReadFromQueue(queue)
				eventChannel <- value
			}
		}()
}


