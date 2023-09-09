package forwarder

import (
	"fmt"
	"time"
	"encoding/json"
	"bytes"
	"net/http"
	"log"
)


type Forwarder struct {
	Webhooks  []string
}

func NewForwarder(webhooks []string) (*Forwarder, error) {
	fs := &Forwarder{
		Webhooks: webhooks,
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


func (f *Forwarder) Send(username, text string) error {
	log.Printf("Sending new message from %s", username)
	
	// Get current time in UTC time zone
	currentTime := time.Now().UTC().Format("2006/01/02 15:04:05") + " ET"

	e := Embed{
	        Author: EmbedAuthor{
			Name: username,
			IconUrl: "https://www.stealth-net.co.uk/EVE/goons_mumble/images/goon.png",
		},
		Description: text,
		Footer: EmbedFooter{
			Text: currentTime,
		},
	}
	
	wm := WebhookMessage{
		Username: "GoonBot",
		Embeds: []Embed{e, },
	}

	postBody, err := json.Marshal(wm)
	if err != nil {
		fmt.Println(err)
	}

	for _, webhook := range f.Webhooks {

		requestBody := bytes.NewBuffer(postBody)

		resp, err := http.Post(
			webhook,
			"application/json",
			requestBody,
		)
		defer resp.Body.Close()
	
		if err != nil {
			log.Printf("Error forwarding to webhook %s: %s", webhook, err)
			continue
		}

		log.Printf("Message forwarded to %s", webhook)
	}

	return nil
}

