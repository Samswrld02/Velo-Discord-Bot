package bot

import (
	"encoding/json"
	"testing"

	"github.com/Samswrld02/Velo-Discord-Bot/types"
	"github.com/bwmarrin/discordgo"
)

//payload exactly as PushTicketCreatedToDiscord in the velo backend pushes it
const ticketCreatedPayload = `{
	"type": "ticket.created",
	"channel_id": "111111111111111111",
	"ticket_id": "9d1c7c1e-0000-4000-8000-000000000001",
	"ticket_number": "42",
	"subject": "Printer werkt niet",
	"priority": "high",
	"channel_type": "email",
	"inbox": "Support",
	"requester": "Jan Jansen",
	"created_by": null,
	"url": "http://frontend.test/tickets/9d1c7c1e-0000-4000-8000-000000000001",
	"occurred_at": "2026-10-19T12:00:00+00:00"
}`

func parseTicket(t *testing.T, payload string) types.Ticket {
	t.Helper()

	var ticket types.Ticket
	if err := json.Unmarshal([]byte(payload), &ticket); err != nil {
		t.Fatalf("could not unmarshal payload: %v", err)
	}

	return ticket
}

func fieldValue(embed *discordgo.MessageEmbed, name string) (string, bool) {
	for _, field := range embed.Fields {
		if field.Name == name {
			return field.Value, true
		}
	}

	return "", false
}

func TestTicketCreatedPayloadUnmarshals(t *testing.T) {
	ticket := parseTicket(t, ticketCreatedPayload)

	if ticket.Type != types.TicketCreated {
		t.Errorf("type = %q, want %q", ticket.Type, types.TicketCreated)
	}
	if ticket.Channel_id != "111111111111111111" {
		t.Errorf("channel_id = %q", ticket.Channel_id)
	}
	if ticket.Inbox != "Support" || ticket.Requester != "Jan Jansen" || ticket.Priority != "high" || ticket.Channel_type != "email" {
		t.Errorf("unexpected ticket fields: %+v", ticket)
	}
	if ticket.Created_by != nil {
		t.Errorf("created_by = %q, want nil", *ticket.Created_by)
	}
}

func TestCreateTicketCreatedMessage(t *testing.T) {
	message := createTicketCreatedMessage(parseTicket(t, ticketCreatedPayload))

	if len(message.Embeds) != 1 {
		t.Fatalf("embeds = %d, want 1", len(message.Embeds))
	}
	embed := message.Embeds[0]

	if embed.Title != "New ticket #42" {
		t.Errorf("title = %q", embed.Title)
	}
	if embed.Description != "Printer werkt niet" {
		t.Errorf("description = %q", embed.Description)
	}
	for name, want := range map[string]string{"Inbox": "Support", "From": "Jan Jansen", "Priority": "high"} {
		if got, ok := fieldValue(embed, name); !ok || got != want {
			t.Errorf("field %q = %q, want %q", name, got, want)
		}
	}
	//tickets from e-mail have no creator, the field is left out
	if _, ok := fieldValue(embed, "Created by"); ok {
		t.Errorf("created by field should be missing when created_by is null")
	}

	button := message.Components[0].(discordgo.ActionsRow).Components[0].(discordgo.Button)
	if button.URL != "http://frontend.test/tickets/9d1c7c1e-0000-4000-8000-000000000001" {
		t.Errorf("button url = %q", button.URL)
	}
}

func TestCreateTicketCreatedMessageShowsCreator(t *testing.T) {
	ticket := parseTicket(t, ticketCreatedPayload)
	creator := "Sam"
	ticket.Created_by = &creator

	got, ok := fieldValue(createTicketCreatedMessage(ticket).Embeds[0], "Created by")
	if !ok || got != "Sam" {
		t.Errorf("created by field = %q, want %q", got, "Sam")
	}
}
