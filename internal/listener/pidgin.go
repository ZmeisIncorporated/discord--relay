package listener

import (
	"fmt"
	"os"
	"time"
	"regexp"
	"log"

	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
)

type PidginListener struct {
	f *forwarder.Forwarder
	logs string
}


func NewPidginListener(forwarder *forwarder.Forwarder, logs string) (PidginListener, error) {
	l := PidginListener{
		f: forwarder,
		logs: logs,
	}
	return l, nil
}


func (l *PidginListener) Send(username string, msg string) {
	err := l.f.Send(username, msg)
	if err != nil {
		log.Printf("Error while sending pidgin message from %s: %s", username, err)
	}
}


type Message struct {
	message string
	username string
	evetime time.Time
}


func parseMessage(text string) (*Message, error) {
	re := regexp.MustCompile(`(?s:\((?P<time>[\d:]+)\) directorbot: (?P<message>.*?)~~~ .*? from (?P<username>.*?) to .*? at (?P<evetime>.*?) EVE ~~~)`)
	parts := re.FindStringSubmatch(text)

	layout := "2006-01-02 15:04:05"
	e := parts[4]
	
	evetime, err := time.Parse(layout, e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing date %s: %w", e, err)
	}
	
	return &Message{
		message: parts[2],
		username: parts[3],
		evetime: evetime,
	}, nil
}

func getMessages(text []byte) []*Message {
	var messages []*Message
	re := regexp.MustCompile(`(?s:\([\d:]+\) directorbot:.*?~~~ .*? ~~~)`)
	m := re.FindAll(text, -1)

	for _, message := range m {
		new_message, err := parseMessage(string(message))
		if err != nil {
			log.Println(err)
			continue
		}
		messages = append(messages, new_message)
	}
	return messages
}


func filterByDate(now time.Time, evetime time.Time) bool {
	diff := now.Sub(evetime)

	// ToDo: looks dangerous
	if diff > 40 * time.Second {
		return false
	}
	
	return true
}


func getMessagesFromFiles(path string) ([]*Message, error) {
	var messages []*Message

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("Error while getting pidgin files from %s: %w", path, err)
	}

	now := time.Now()

	for _, filename := range files {
		if filename.IsDir() {
			continue
		}

		filepath := fmt.Sprintf("%s/%s", path, filename.Name())
		rawdata, err := os.ReadFile(filepath)
		if err != nil {
			log.Printf("Error while reading file %s: %s", filepath, err)
			continue
		}

		for _, m := range getMessages(rawdata) {
			if filterByDate(now, m.evetime) {
				messages = append(messages, m)
			}
		}
	}

	log.Printf("Got %d new messages from pidgin", len(messages))
	return messages, nil
}


func (l *PidginListener) Run() {
	log.Printf("Pidgin listener started")
	for {
		select {
		case <- time.After(30 * time.Second):

			messages, err := getMessagesFromFiles(l.logs)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, m := range messages {
				l.Send(m.username, m.message)
			}
		}
	}

}

