package listener

import (
	"fmt"
	"os"
	"time"
	"regexp"

	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
)

type PidginListener struct {
	f *forwarder.Forwarder
	ch string
	logs string
	delete bool
}


func NewPidginListener(forwarder *forwarder.Forwarder, ch string, logs string, delete bool) (PidginListener, error) {
	l := PidginListener{
		f: forwarder,
		ch: ch,
		logs: logs,
		delete: delete,
	}
	return l, nil
}


func (l *PidginListener) Send(username string, msg string) {
	err := l.f.Send(username, msg, l.ch)
	if err != nil {
		fmt.Printf("Error while sending pidgin log file: %s", err)
	}
}


type Message struct {
	message string
	username string
	evetime time.Time
}


func parseMessage(text string) *Message {
	re := regexp.MustCompile(`(?s:\((?P<time>[\d:]+)\) directorbot: (?P<message>.*?)~~~ .*? from (?P<username>.*?) to .*? at (?P<evetime>.*?) EVE ~~~)`)
	parts := re.FindStringSubmatch(text)

	layout := "2006-01-02 15:04:05"
	e := parts[4]
	
	evetime, err := time.Parse(layout, e)
	if err != nil {
		fmt.Println("Error while parsing date", err)
	}
	
	return &Message{
		message: parts[2],
		username: parts[3],
		evetime: evetime,
	}
}

func getMessages(text []byte) []*Message {
	var messages []*Message
	re := regexp.MustCompile(`(?s:\([\d:]+\) directorbot:.*?~~~ .*? ~~~)`)
	m := re.FindAll(text, -1)
	for _, message := range m {
		messages = append(messages, parseMessage(string(message)))
	}
	return messages
}


func filterByDate(now time.Time, evetime time.Time) bool {
	diff := now.Sub(evetime)

	// ToDo: looks dangerous
	if diff > 30 * time.Second {
		return false
	}
	
	return true
}


func getMessagesFromFiles(path string) []*Message {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error while getting files")
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
			fmt.Println("Error while reading file")
		}

		for _, m := range getMessages(rawdata) {
			if filterByDate(now, m.evetime) {
				messages = append(messages, m)
			}
		}
	}

	return messages
}


func (l *PidginListener) Run() {
	for {
		select {
		case <- time.After(5 * time.Second):
	
			messages := getMessagesFromFiles(l.logs)

			for _, m := range messages {
				l.Send(m.username, m.message)
			}
		}
	}

}

