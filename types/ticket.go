package types

//event types the velo backend pushes onto the queue
const (
	TicketAssigned = "ticket.assigned"
	TicketCreated = "ticket.created"
)

type Ticket struct {
	Type string `json:"type"`
	Discord_id string `json:"discord_user_id"`
	Ticket_it string `json:"ticket_id"`
	Ticket_number string `json:"ticket_number"`
	Subject string `json:"subject"`
	Assigned_by string `json:"assigned_by"`
	Url string `json:"url"`
	Occured_at string `json:"occured_at"`

	//ticket.created: channel of the linked guild the message goes to
	Channel_id string `json:"channel_id"`
	Priority string `json:"priority"`
	Channel_type string `json:"channel_type"`
	Inbox string `json:"inbox"`
	Requester string `json:"requester"`
	//nil for tickets that came in by e-mail
	Created_by *string `json:"created_by"`
}