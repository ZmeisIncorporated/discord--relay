package listener

import (
	"fmt"
	"io/ioutil"
	"os"
	"golang.org/x/net/html"
	"strings"
	"errors"
	"bytes"
	"io"

	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
)

type PidginListener struct {
	f *forwarder.Forwarder
	ch string
	logs string
}

func NewPidginListener(forwarder *forwarder.Forwarder, ch string, logs string) (PidginListener, error) {
	l := PidginListener{
		f: forwarder,
		ch: ch,
		logs: logs,
	}
	return l, nil
}


func (l *PidginListener) Send(msg string) {
	err := l.f.Send("[pidgin]", msg, l.ch)
	if err != nil {
		fmt.Printf("Error while sending pidgin log file: %s", err)
	}
}

func Body(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)

	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}


func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func (l *PidginListener) MsgParser(text string) string {
	raw_message, err := html.Parse(strings.NewReader(text))
	if err != nil {
		l.Send("Error while parsing pidgin message")
	}

	bn, err := Body(raw_message)
	if err != nil {
		l.Send("Error while parsing pidgin message")
	}

	body := renderNode(bn)
	return body
}


func (l *PidginListener) Run() {
	files, err := ioutil.ReadDir(l.logs)
	if err != nil {
		l.Send("Error while reading pidgin log files")
		return
	}

	for _, file := range files {
		f, err := os.ReadFile(fmt.Sprintf("%s/%s", l.logs, file.Name()))
		if err != nil {
			l.Send(fmt.Sprintf("Error while reading logfile %s: %s", file.Name(), err))
		}
		fs := string(f)
		body := l.MsgParser(fs)
		l.Send(body)
	}
}

