package listener

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
)

const (
	pidgin        = "pidgin"
	checkTimeout  = 30 * time.Second
	filterTimeout = checkTimeout + 90*time.Second
)

type PidginListener struct {
	f    *forwarder.Forwarder
	logs string
}

func NewPidginListener(forwarder *forwarder.Forwarder, logs string) PidginListener {
	l := PidginListener{
		f:    forwarder,
		logs: logs,
	}
	return l
}

type Message struct {
	message  string
	username string
	evetime  time.Time
}

func parseMessage(text string) (*Message, error) {
	re := regexp.MustCompile(`(?s:\((?P<time>[\d:]+)\) directorbot: (?P<message>.*?)~~~ .*? from (?P<username>.*?) to .*? at (?P<evetime>.*?) EVE ~~~)`)
	parts := re.FindStringSubmatch(text)

	layout := "2006-01-02 15:04:05"
	e := parts[4]

	evetime, err := time.Parse(layout, e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing date: %w", err)
	}

	return &Message{
		message:  parts[2],
		username: parts[3],
		evetime:  evetime,
	}, nil
}

func getMessages(text []byte) ([]*Message, error) {
	var messages []*Message
	re := regexp.MustCompile(`(?s:\([\d:]+\) directorbot:.*?~~~ .*? ~~~)`)
	m := re.FindAll(text, -1)

	for _, message := range m {
		new_message, err := parseMessage(string(message))
		if err != nil {
			return nil, err
		}
		messages = append(messages, new_message)
	}
	return messages, nil
}

func filterByDate(now time.Time, evetime time.Time) bool {
	diff := now.Sub(evetime)

	if diff > filterTimeout {
		return false
	}

	return true
}

func getMessagesFromFiles(path string) ([]*Message, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("Error while getting files: %w", err)
	}

	var messages []*Message
	now := time.Now()

	for _, filename := range files {
		if filename.IsDir() {
			continue
		}

		filepath := fmt.Sprintf("%s/%s", path, filename.Name())
		rawdata, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("Error while reading file: %w", err)
		}

		new_messages, err := getMessages(rawdata)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing file: %w", err)
		}

		for _, m := range new_messages {
			if filterByDate(now, m.evetime) {
				messages = append(messages, m)
			}
		}
	}

	return messages, nil
}

func (l *PidginListener) Start() {
	log.Println("Starting pidgin listener")
	go func() {
		var lastUsername, lastMessage string = "", ""
		for {
			select {
			case <-time.After(checkTimeout):

				messages, err := getMessagesFromFiles(l.logs)
				if err != nil {
					msg := fmt.Sprintf("PidginListener error: %s", err)
					log.Println(msg)
					l.f.AdmSend(pidgin, msg)
				}

				if len(messages) > 0 {
					log.Printf("Got %d last messages from pidgin", len(messages))
				}

				for _, m := range messages {
					if lastUsername == m.username && lastMessage == m.message {
						log.Printf("Message from %s was send already, skip", m.username)
						continue
					}
					l.f.WebSend(m.username, m.message)
					lastUsername, lastMessage = m.username, m.message
				}
			}
		}
	}()
}
