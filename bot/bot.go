package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/Samswrld02/Velo-Discord-Bot/types"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	redisconnection "github.com/Samswrld02/Velo-Discord-Bot/redis-connection"
	veloapi "github.com/Samswrld02/Velo-Discord-Bot/velo-api"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
	wg sync.WaitGroup
	logger *logrus.Logger
	config *types.Config
	velo *veloapi.Client
}

func (b *Bot) Run() {
	//slash command handler, added before Open so no interaction is missed
	b.session.AddHandler(b.handleInteraction)

	//open session
	if err := b.session.Open(); err != nil {
		b.logger.Fatal("Could not open Discord session " + err.Error())
	}

	defer b.session.Close()

	b.registerCommands()

	b.printAscii()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//event handler
	b.handleQueue(ctx)

	b.logger.Info("bot is running")

	b.handleGracefulshutdown(cancel)
}

func (b *Bot) handleGracefulshutdown(cancel context.CancelFunc) {
	//gracefull shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	
	<-c
	b.logger.Info("Bot is shutting down")
	//signal context to close
	cancel()

	//wait for goroutines to finish
	b.wg.Wait()
	fmt.Println("all events have been processed, bot has been closed! Goodbye")
}


func NewBot(logger *logrus.Logger, config *types.Config) *Bot{
	//create a discord session
	discord, err := discordgo.New("Bot " +  config.DiscordBotToken)
	if err != nil {
		panic(err)
	}


	velo := veloapi.NewClient(config.VeloAPIURL, config.VeloBotToken)

	return &Bot{session: discord, wg: sync.WaitGroup{}, logger: logger, config: config, velo: velo}
}

//method for checking if ticket event was put into the queue
func (b *Bot) handleQueue(ctx context.Context) {

	//retrieve redis db for the queue
	db, err := strconv.Atoi(b.config.RedisDb)
	if err != nil {
		panic("error redis db env variable " + err.Error())
	}
	
	//create redis connection
	r := redisconnection.NewRedisClient(b.config.RedisAddr, b.config.RedisPw, db)

	//channel for the queued events
	eventChannel := make(chan []byte)

	//worker that reads data from the redis queue into the channel
	b.readEventWorker(eventChannel, r, ctx)

	//specified ammount of workers that process data from the channel concurrently
	b.handleTicketWorkers(5, eventChannel)
}

//method for processing data into struct
func (b *Bot) handleData(data []byte) {
	//import ticket struct
	var ticket types.Ticket;

	//unmarshal data into struct
	if err := json.Unmarshal(data, &ticket); err != nil  {
		b.logger.Error("error unmarshaling data: " + err.Error())

		//early return if data is malformed
		return 
	}

	ch, err := b.session.UserChannelCreate(ticket.Discord_id)

	if err != nil {
		b.logger.Error("Error occurred: " + err.Error())
		return
	}

	message := createMessage(ticket)

	_, err = b.session.ChannelMessageSendComplex(ch.ID, message)

	if err != nil {
		b.logger.Error(err.Error())
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
		//add to wg for every worker
		b.wg.Add(1)

		//dispatch worker
		go func(){
			defer b.wg.Done()
			//read from channel, blocks if channel is empty
			for job := range eventChannel {
				b.handleData(job)
			}

			//when eventchannel is closed go routine will finish on its own and return
		}()
	}
}

//processo/worker of redis queue
func (b *Bot) readEventWorker(eventChannel chan []byte, r *redisconnection.RedisC, ctx context.Context) {
		queue := b.config.RedisQueue
		
		go func() {		
		//loop infinitely blpop will block if queue is empty
			for {
				select{
				case <-ctx.Done():
					//handle stopping, close eventchannel
					close(eventChannel)
					return
				default:
					value, err := r.ReadFromQueue(ctx, queue)
					if err != nil{
						if _, ok := errors.AsType[*net.OpError](err); ok {
							b.logger.Panic(err.Error())
						}
						if !errors.Is(err, redis.Nil) {
							b.logger.Error(err.Error())
						}
						continue
					}
					eventChannel <- value
				}
			}
		}()
}

func (b *Bot) printAscii() {
	fmt.Print(`
█   █ █████ █      ███     ████   ███  █████    █     ███ █   █ █████ 
█   █ █     █     █   █    █   █ █   █   █      █      █  █   █ █     
█   █ ████  █     █   █    ████  █   █   █      █      █  █   █ ████  
 █ █  █     █     █   █    █   █ █   █   █      █      █   █ █  █     
  █   █████ █████  ███     ████   ███    █      █████ ███   █   █████ `)
}


