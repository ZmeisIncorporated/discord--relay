package forwarder

import (
	"fmt"
	"log"
	"time"
	"encoding/json"
	"bytes"
	"net/http"
)


// Forwarder Wrapper for the discord session
type Forwarder struct {
	AdmHooks   []string
	WebHooks   []string
	WebIn      chan *WebhookMessage
	AdmIn      chan *WebhookMessage
}

// NewForwarder takes in a token and returns a Forward Session
func NewForwarder(webhooks, admhooks []string) (*Forwarder, error) {
	fs := &Forwarder{
		AdmHooks: admhooks, 
		WebHooks: webhooks,
		WebIn:    make(chan *WebhookMessage),
		AdmIn:    make(chan *WebhookMessage),
	}
	return fs, nil
}

type EmbedAuthor struct {
	Name      string  `json:"name"`
	IconUrl   string  `json:"icon_url"`
}

type EmbedFooter struct {
	Text      string  `json:"text"`
}

type Embed struct {
	Author       EmbedAuthor  `json:"author"`
	Description  string       `json:"description"`
	Footer       EmbedFooter  `json:"footer"`
}

type WebhookMessage struct {
	Username  string  `json:"username"`
	Embeds    []Embed `json:"embeds"`
}


func (f *Forwarder) Send2Hooks(wm *WebhookMessage, webhooks []string) error {
	postBody, err := json.Marshal(wm)
	if err != nil {
		fmt.Println(err)
	}
	requestBody := bytes.NewBuffer(postBody)

	for _, webhook := range webhooks {

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


func CreateWebhookMessage(username, message string) *WebhookMessage {
	// Get current time in UTC time zone
	currentTime := time.Now().UTC().Format("2006/01/02 15:04:05") + " ET"

	e := Embed{
	        Author: EmbedAuthor{
			Name: username,
			IconUrl: "https://www.stealth-net.co.uk/EVE/goons_mumble/images/goon.png",
		},
		Description: message,
		Footer: EmbedFooter{
			Text: currentTime,
		},
	}
	
	return &WebhookMessage{
		Username: "GoonBot",
		Embeds: []Embed{e, },
	}
}


func (f *Forwarder) WebSend(username, message string) {
	wm := CreateWebhookMessage(username, message)
	f.WebIn <- wm
}

func (f *Forwarder) AdmSend(username, message string) {
	wm := CreateWebhookMessage(username, message)
	f.AdmIn <- wm
}


func (f *Forwarder) Start() {
	go func() {
		for {
			select {
			case w := <- f.WebIn:
				err := f.Send2Hooks(w, f.WebHooks)
				if err != nil {
					log.Printf("%s", err)
				}
			case a := <- f.AdmIn:
				err := f.Send2Hooks(a, f.AdmHooks)
				if err != nil {
					log.Printf("%s", err)
				}
			}
		}
	}()

}
