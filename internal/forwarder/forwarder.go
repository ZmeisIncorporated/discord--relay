package forwarder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ZmeisIncorporated/discord--relay/internal/config"
)

// Message contains discord message and list of hooks to send to
type Message struct {
	Message  *WebhookMessage
	Webhooks []string
}

type Forwarder struct {
	AdmHooks []string
	WebHooks []string
	Filters  map[string]config.Filter
	Sender   chan *Message
	IconUrl  string
	Botname  string
}

// NewForwarder takes in a token and returns a Forward Session
func NewForwarder(webhooks, admhooks []string, filters map[string]config.Filter, iconUrl, botname string) (*Forwarder, error) {
	fs := &Forwarder{
		AdmHooks: admhooks,
		WebHooks: webhooks,
		Filters:  filters,
		Sender:   make(chan *Message),
		IconUrl:  iconUrl,
		Botname:  botname,
	}
	return fs, nil
}

type EmbedAuthor struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type Embed struct {
	Author      EmbedAuthor `json:"author"`
	Description string      `json:"description"`
	Footer      EmbedFooter `json:"footer"`
}

type WebhookMessage struct {
	Username string  `json:"username"`
	Embeds   []Embed `json:"embeds"`
}

// Send2Hooks sends message to discord hooks
func (f *Forwarder) Send2Hooks(wm *WebhookMessage, webhooks []string) error {
	postBody, err := json.Marshal(wm)
	if err != nil {
		fmt.Println(err)
	}

	for _, webhook := range webhooks {
		requestBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(
			webhook,
			"application/json",
			requestBody,
		)
		defer resp.Body.Close()

		if err != nil {
			return fmt.Errorf("error forwarding to webhook %s", err)
		}

		log.Printf("Message forwarded to %s", webhook)
	}

	return nil
}

func CreateWebhookMessage(username, message, iconUrl, botname string) *WebhookMessage {
	// Get current time in UTC time zone
	currentTime := time.Now().UTC().Format("2006/01/02 15:04:05") + " ET"

	e := Embed{
		Author: EmbedAuthor{
			Name:    username,
			IconUrl: iconUrl,
		},
		Description: message,
		Footer: EmbedFooter{
			Text: currentTime,
		},
	}

	return &WebhookMessage{
		Username: botname,
		Embeds:   []Embed{e},
	}
}

// RouteMessage send message to destination depends on filters settings
func (f *Forwarder) RouteMessage(username, message string) {
	wm := CreateWebhookMessage(username, message, f.IconUrl, f.Botname)

	for filterName, filter := range f.Filters {
		for _, pattern := range filter.Patterns {

			if strings.Contains(message, pattern) {
				log.Printf("Filter %s worked", filterName)
				f.Sender <- &Message{
					wm,
					filter.Webhooks,
				}
				if !filter.Propagate {
					log.Printf("Filter %s propagate is off, message will not moved to main webhooks", filterName)
					return
				}
				log.Printf("Filter %s propagate is on, message will moved to main webhooks", filterName)
			}
		}
	}

	log.Printf("Process message to default webhooks")
	f.Sender <- &Message{
		wm,
		f.WebHooks,
	}
}

// WebSend sends messages to main webhooks without filtering
func (f *Forwarder) WebSend(username, message string) {
	wm := CreateWebhookMessage(username, message, f.IconUrl, f.Botname)
	f.Sender <- &Message{
		wm,
		f.WebHooks,
	}
}

// AdmSend sends messages to admin webhooks without filtering
func (f *Forwarder) AdmSend(username, message string) {
	wm := CreateWebhookMessage(username, message, f.IconUrl, f.Botname)
	f.Sender <- &Message{
		wm,
		f.AdmHooks,
	}
}

func (f *Forwarder) Start() {
	go func() {
		for {
			select {
			case s := <-f.Sender:
				err := f.Send2Hooks(s.Message, s.Webhooks)
				if err != nil {
					log.Printf("%s", err)
				}
			}
		}
	}()

}
