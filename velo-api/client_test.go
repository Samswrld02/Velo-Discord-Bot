package veloapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var linkRequest = LinkGuildRequest{
	GuildID:       "111111111111111111",
	GuildName:     "Velo Team",
	ChannelID:     "222222222222222222",
	ChannelName:   "velo-meldingen",
	DiscordUserID: "999999999999999999",
}

//fake velo api that answers with the given status and message
func fakeVelo(t *testing.T, status int, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/discord/guilds" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("X-Velo-Bot-Token"); got != "secret-bot-token" {
			t.Errorf("bot token header = %q, want %q", got, "secret-bot-token")
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("accept header = %q, want application/json", got)
		}

		var body LinkGuildRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("could not decode body: %v", err)
		}
		if body != linkRequest {
			t.Errorf("body = %+v, want %+v", body, linkRequest)
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
	}))
}

func TestLinkGuildReturnsTheMessageOnSuccess(t *testing.T) {
	server := fakeVelo(t, http.StatusCreated, "Kanaal #velo-meldingen gekoppeld aan Velo ✅")
	defer server.Close()

	message, err := NewClient(server.URL, "secret-bot-token").LinkGuild(context.Background(), linkRequest)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message != "Kanaal #velo-meldingen gekoppeld aan Velo ✅" {
		t.Errorf("message = %q", message)
	}
}

func TestLinkGuildReturnsVelosMessageForAKnownProblem(t *testing.T) {
	server := fakeVelo(t, http.StatusConflict, "Je zit in meerdere organisaties; koppelen via Discord kan dan nog niet.")
	defer server.Close()

	message, err := NewClient(server.URL, "secret-bot-token").LinkGuild(context.Background(), linkRequest)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message != "Je zit in meerdere organisaties; koppelen via Discord kan dan nog niet." {
		t.Errorf("message = %q", message)
	}
}

func TestLinkGuildReturnsAnErrorWhenTheBotIsRejected(t *testing.T) {
	server := fakeVelo(t, http.StatusUnauthorized, "Invalid bot token.")
	defer server.Close()

	_, err := NewClient(server.URL, "secret-bot-token").LinkGuild(context.Background(), linkRequest)

	if err == nil {
		t.Fatal("expected an error for a 401, got nil")
	}
}

func TestLinkGuildReturnsAnErrorWhenVeloIsDown(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	url := server.URL
	server.Close()

	_, err := NewClient(url, "secret-bot-token").LinkGuild(context.Background(), linkRequest)

	if err == nil {
		t.Fatal("expected an error when velo is unreachable, got nil")
	}
}
