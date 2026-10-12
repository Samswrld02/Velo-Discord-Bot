package types

type Ticket struct {
	Id string `json:"ticket_id"`
	Number string `json:"ticket_number"`
	Url string `json:"ticket_url"`
	Assignee string `json:"ticket_assignee"`
	Discord_id string `json:"discord_id"`
}