package types

type Ticket struct {
	Type string `json:"type"`
	Discord_id string `json:"discord_user_id"`
	Ticket_it string `json:"ticket_id"`
	Ticket_number string `json:"ticket_number"`
	Subject string `json:"subject"`
	Assigned_by string `json:"assigned_by"`
	Url string `json:"url"`
	Occured_at string `json:"occured_at"`
	
}