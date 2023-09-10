package monitor

import (
	"log"
	"time"

	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
	gops "github.com/mitchellh/go-ps"
)

const finch = "finch"

type Monitor struct {
	f *forwarder.Forwarder
}

func searchFinch() bool {
	processes, err := gops.Processes()
	if err != nil {
		log.Printf("Unable to get running processes: %s", err)
		return false
	}

	for _, p := range processes {
		if p.Executable() == "finch" {
			return true
		}
	}

	return false
}

func (m *Monitor) Start() {
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				if f := searchFinch(); !f {
					msg := "Finch process not found, need to restart it manually"
					m.f.AdmSend("[monitor]", msg)
					log.Println(msg)
				}
			}
		}
	}()
}

func NewMonitor(f *forwarder.Forwarder) *Monitor {
	return &Monitor{
		f: f,
	}
}
