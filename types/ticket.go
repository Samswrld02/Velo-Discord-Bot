package types

//event types the velo backend pushes onto the queue
const (
	TicketAssigned = "ticket.assigned"
	TicketCreated = "ticket.created"
)

//color values
const (
	ColorLow      = 0x3498db // Blue (#3498db) - Information / Informational
	ColorNormal   = 0x2ecc71 // Green (#2ecc71) - Healthy / Normal operation
	ColorHigh     = 0xe67e22 // Orange (#e67e22) - Warning / High priority
	ColorCritical = 0xe74c3c // Red (#e74c3c) - Alert / Critical error
)

var Priories = map[string]int{
	"low" : ColorLow,
	"normal": ColorNormal,
	"high": ColorHigh,
	"critical": ColorCritical,
}

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