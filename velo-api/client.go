package veloapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//client for calling the velo backend as the bot
type Client struct {
	baseURL string
	token string
	http *http.Client
}

//what the bot sends when someone runs /kanaal
type LinkGuildRequest struct {
	GuildID string `json:"guild_id"`
	GuildName string `json:"guild_name"`
	ChannelID string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	DiscordUserID string `json:"discord_user_id"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token: token,
		//discord gives us 15 minutes after a deferred reply, but nobody wants to wait that long
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

//links the guild + channel to velo, returns the dutch message velo wants to show the user
func (c *Client) LinkGuild(ctx context.Context, req LinkGuildRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/discord/guilds", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	//without this laravel answers validation errors with a redirect instead of json
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Velo-Bot-Token", c.token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res messageResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("velo api returned status %d with an unreadable body: %w", resp.StatusCode, err)
	}

	switch resp.StatusCode {
	//velo wrote these messages for the user (linked, not linked, not admin, conflict)
	case http.StatusOK, http.StatusCreated, http.StatusNotFound, http.StatusForbidden, http.StatusConflict:
		return res.Message, nil
	//everything else (wrong bot token, validation, server error) is our problem, not the user's
	default:
		return "", fmt.Errorf("velo api returned status %d: %s", resp.StatusCode, res.Message)
	}
}
