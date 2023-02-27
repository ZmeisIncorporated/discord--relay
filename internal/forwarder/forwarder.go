package forwarder

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Forwarder Wrapper for the discord session
type Forwarder struct {
	*discordgo.Session
	IsWebHook bool
	Channels  map[string]string
	log       *logrus.Logger
}

// NewForwarder takes in a token and returns a Forward Session
func NewForwarder(token string, webhook bool, chans map[string]string, log *logrus.Logger) (*Forwarder, error) {
	d, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, err
	}
	d.Identify.Presence.Status = string(discordgo.StatusDoNotDisturb)
	err = d.Open()
	if err != nil {
		return nil, err
	}
	fs := &Forwarder{d, webhook, chans, log}
	return fs, nil
}

// Send send a message here and we'll figure out how
func (f *Forwarder) Send(username, text, channelID string) error {
	if f.IsWebHook {
		return f.ToWebhook(username, text, channelID)
	}
	return f.ToMessage(username, text, channelID)
}

// ToMessage forwards a message to the specific chan as a message
func (f *Forwarder) ToMessage(username, text string, channelID string) error {
	_, err := f.ChannelMessageSend(channelID, fmt.Sprintf("**%s**: %s", username, text))
	if err != nil {
		return fmt.Errorf("error forwarding to Message %s", err)
	}
	return nil
}

// ToWebhook forwards a message to the specific chan as a webhook
func (f *Forwarder) ToWebhook(username, text string, channelID string) error {
	ws, err := f.ChannelWebhooks(channelID)
	if err != nil {
		return fmt.Errorf("error looking up webhook %s", err)
	}
	var wh *discordgo.Webhook
	if len(ws) == 0 {
		wh, err = f.WebhookCreate(channelID, "botman", "")
		if err != nil {
			return fmt.Errorf("error creating webhook %s", err)
		}
	} else {
		wh = ws[0]
	}

	// Get current time in UTC time zone
	currentTime := time.Now().UTC().Format("2006/01/02 15:04:05") + " ET"

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "",
			IconURL: "",
		},
		Title:       "",
		Description: text,
		Footer: &discordgo.MessageEmbedFooter{
			Text: currentTime,
		},
	}

	_, err = f.WebhookExecute(
		wh.ID,
		wh.Token,
		true,
		&discordgo.WebhookParams{
			Content:         "",
			Username:        "",
			AvatarURL:       "",
			TTS:             false,
			Embeds:          []*discordgo.MessageEmbed{embed},
			AllowedMentions: nil,
		})

	if err != nil {
		return fmt.Errorf("error forwarding to webhook %s", err)
	}
	return nil
}

// MessageCreate discord handler for new messages
func (f *Forwarder) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var forwardChanID string
	var ok bool
	// Are we listening to this channel?
	if forwardChanID, ok = f.Channels[m.ChannelID]; !ok {
		// f.log.Debugf("didnt find channel: %s", m.ChannelID)
		// f.log.Debugf("all channels: %v", f.Channels)
		return
	}

	if m.Content == "" {
		l, err := s.ChannelMessages(m.ChannelID, 1, "", "", m.ID)
		if err != nil {
			f.log.Errorf("unable to get messages: %s", err)
			return
		}
		m.Content = l[0].Content
		m.Attachments = l[0].Attachments
		m.Author = l[0].Author
	}

	// Get sender info
	srcMember, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		f.log.Debugf("didnt find srcMember: %s", m.Author.Username)
		return
	}

	// Prevent this from relaying @mentions through
	// Todo previously we converted IDs to nice names
	noAtsReg := regexp.MustCompile(`@(\S+)`)
	m.Content = noAtsReg.ReplaceAllString(m.Content, "**@ $1**")

	// Pull out the links and make them text again
	var links string
	if m.Attachments != nil {
		for _, a := range m.Attachments {
			links += a.URL + " "
		}
	}
	// convert username into a legible format
	var username string
	if srcMember.Nick != "" {
		username = fmt.Sprintf("%s (%s#%s)", srcMember.Nick, srcMember.User.Username, srcMember.User.Discriminator)
	} else {
		username = fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator)
	}

	// Finally send a nicely formated message
	// f.log.Debugf("u: %s - msg: %s - chanid: %d\n", username, fmt.Sprintf("%s %s", m.Content, links), forwardChanID)
	err = f.Send(username, fmt.Sprintf("%s %s", m.Content, links), forwardChanID)
	if err != nil {
		f.log.Errorf("error sending: %s", err)
	}

}
