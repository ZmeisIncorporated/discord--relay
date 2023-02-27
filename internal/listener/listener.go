package listener

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/justindh/discord-relay/internal/forwarder"
)

//Listeners all of the incloming updates.
type Listeners struct {
	Sessions []*discordgo.Session
	log      *logrus.Logger
}

// NewListeners bulk creates all listeners and gets them running
func NewListeners(tokens []string, forwarder *forwarder.Forwarder, log *logrus.Logger) (Listeners, error) {
	l := Listeners{Sessions: make([]*discordgo.Session, 0), log: log}
	for _, t := range tokens {
		sess, err := newListener(t, forwarder)
		if err != nil {
			return l, err
		}
		l.Sessions = append(l.Sessions, sess)
	}
	return l, nil
}

func newListener(token string, forwarder *forwarder.Forwarder) (*discordgo.Session, error) {
	sess, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	// We want to appear offline
	sess.Identify.Presence.Status = string(discordgo.StatusInvisible)
	// We want to filter to only messages (less clutter)
	sess.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
	// Add a message processor
	sess.AddHandler(forwarder.MessageCreate)
	// Change the defaults so we look less like a bot
	sess.Identify.Properties.Browser = "Chrome"
	sess.Identify.Properties.Device = ""
	sess.Identify.Properties.OS = "Windows"
	sess.Identify.Properties.Referer = "https://www.google.com/"
	sess.Identify.Properties.ReferringDomain = "www.google.com"
	sess.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
	sess.Debug = false
	// open the session and start listening for messages
	err = sess.Open()
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// Close attempts to stop and clean up all listeners
func (l *Listeners) Close() {
	for _, s := range l.Sessions {
		s.Close()
	}
}
