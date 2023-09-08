package forwarder

import (
	"fmt"
	"regexp"
	"time"
	"encoding/json"
	"bytes"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	noAtsReg = regexp.MustCompile(`@(\S+)`)
)

// Forwarder Wrapper for the discord session
type Forwarder struct {
	webhook   string
	log       *logrus.Logger
}

// NewForwarder takes in a token and returns a Forward Session
func NewForwarder(webhook string, log *logrus.Logger) (*Forwarder, error) {
	fs := &Forwarder{
		webhook: webhook,
		log: log,
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


// Send forwards a message to the specific chan as a webhook
func (f *Forwarder) Send(username, text string) error {
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
	responseBody := bytes.NewBuffer(postBody)

	
	_, err = http.Post(
		"https://discord.com/api/webhooks/1149623927560208405/Bh8S3PxfzAWDNxqmOH30V1dmDLAIS1bNj8n_2XGEao0pG2bdMXGTdmSGr0po8eaMOPGj",
		"application/json",
		responseBody,
	)
	
	
	if err != nil {
		return fmt.Errorf("error forwarding to webhook %s", err)
	}
	return nil
}

