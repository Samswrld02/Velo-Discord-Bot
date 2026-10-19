package bot

import (
	"context"
	"time"

	veloapi "github.com/Samswrld02/Velo-Discord-Bot/velo-api"
	"github.com/bwmarrin/discordgo"
)

const linkChannelCommand = "kanaal"

//only members that can manage channels see the command
var manageChannels int64 = discordgo.PermissionManageChannels

//slash commands the bot registers on startup
var commands = []*discordgo.ApplicationCommand{
	{
		Name: linkChannelCommand,
		Description: "Koppel dit kanaal aan Velo",
		DefaultMemberPermissions: &manageChannels,
	},
}

//registers the slash commands, must run after the session is open (needs the bot's own id)
func (b *Bot) registerCommands() {
	//bulk overwrite replaces all commands, so removed commands disappear too
	//dev guild = instantly available on that server, empty = global
	_, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, b.config.DiscordDevGuildID, commands)
	if err != nil {
		b.logger.Fatal("Could not register slash commands " + err.Error())
	}
}

//handler for every interaction (slash commands, buttons, ...)
func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand || i.ApplicationCommandData().Name != linkChannelCommand {
		return
	}

	//in a DM there is no guild and i.Member is nil
	if i.GuildID == "" {
		b.respond(i, "Gebruik /kanaal in het kanaal van een server waar Velo moet posten.")
		return
	}

	//discord wants an answer within 3 seconds, so first say "busy" and edit the message later
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral},
	})
	if err != nil {
		b.logger.Error("could not defer interaction: " + err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message, err := b.velo.LinkGuild(ctx, veloapi.LinkGuildRequest{
		GuildID: i.GuildID,
		GuildName: b.guildName(i.GuildID),
		ChannelID: i.ChannelID,
		ChannelName: b.channelName(i.ChannelID),
		DiscordUserID: i.Member.User.ID,
	})
	if err != nil {
		b.logger.Error("could not link guild: " + err.Error())
		message = "Velo is nu niet bereikbaar. Probeer het later opnieuw."
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &message}); err != nil {
		b.logger.Error("could not edit interaction response: " + err.Error())
	}
}

//direct ephemeral answer, only the person who ran the command sees it
func (b *Bot) respond(i *discordgo.InteractionCreate, message string) {
	err := b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		b.logger.Error("could not respond to interaction: " + err.Error())
	}
}

//name from the session cache, falls back to the api, and to the id if both fail
func (b *Bot) guildName(guildID string) string {
	if guild, err := b.session.State.Guild(guildID); err == nil {
		return guild.Name
	}
	if guild, err := b.session.Guild(guildID); err == nil {
		return guild.Name
	}
	return guildID
}

func (b *Bot) channelName(channelID string) string {
	if channel, err := b.session.State.Channel(channelID); err == nil {
		return channel.Name
	}
	if channel, err := b.session.Channel(channelID); err == nil {
		return channel.Name
	}
	return channelID
}
