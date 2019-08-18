package yeelight

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// notes:
// max 4 parallel opened TCP connections
// quota: 60 commands per minute (for one device)
// quota: 144 commands per minute for all devices
// TODO: Returns response objects too, not only error
// TODO: Export interface only, not whole struct
type Bulb struct {
	standardCommands
	commonCommands

	// Namespace to control "background" capabilities (device must support it)
	Bg backgroundLightCommands

	Ip   string
	Port int

	conn       net.Conn
	results    map[int]chan Response
	resultsMtx sync.Mutex
}

func (b *Bulb) Connect() error {
	destination := fmt.Sprintf("%s:%d", b.Ip, b.Port)

	conn, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}

	go b.responseProcessor()

	b.conn = conn
	return nil
}

func (b *Bulb) Disconnect() error {
	err := b.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// NewBulb creates Bulb instance, default protocol port: 55443
func NewBulb(ip string, port int) *Bulb {
	bulb := &Bulb{
		standardCommands{},
		commonCommands{},
		backgroundLightCommands{},
		ip,
		port,
		nil,
		make(map[int]chan Response),
		sync.Mutex{},
	}
	// I know It looks badly, but "It is working? It is working"
	bulb.standardCommands.commander = bulb
	bulb.commonCommands.commander = bulb
	bulb.Bg.commander = bulb
	bulb.Bg.prefix = "bg_"
	return bulb
}

func (b *Bulb) executeCommand(c partialCommand) error {
	respChan := make(chan Response)

	// preparing request ID to be able to monitor and wait for response
	b.resultsMtx.Lock()
	id, err := b.findFirstFreeIntKey()
	if err != nil {
		b.resultsMtx.Unlock()
		return err
	}
	b.results[id] = respChan
	b.resultsMtx.Unlock()

	defer func(ch chan Response, id int) {
		close(ch)
		delete(b.results, id)
	}(respChan, id)

	realCommand := newCompleteCommand(c, id)
	message, err := json.Marshal(realCommand)
	if err != nil {
		return err
	}
	log.Printf("[%s] request: %s\n", b.Ip, message)
	message = append(message, CR, LF)

	_, err = b.conn.Write(message)
	if err != nil {
		return err
	}

	// waiting for response on that request
	resp := <-respChan
	return resp.ok()
}

func openSocket(host string, min, max int) (net.Listener, int, error) {
	if min > max {
		return nil, 0, errors.New("min value cannot be greather than max value")
	}
	if min < 0 || max > 65535 {
		return nil, 0, errors.New("range must be 0 - 65535")
	}

	for port := min; port <= max; port++ {
		var ip = "" // binding on all interfaces
		address := fmt.Sprintf("%s:%d", ip, port)

		listener, err := net.Listen("tcp", address)
		if err != nil {
			continue
		}
		return listener, port, nil
	}
	return nil, 0, errors.New("no available free ports in given range")

}

// keysExists returns a bool when given map contains all of given key names
func keysExists(m map[string]interface{}, keys ...string) bool {
	var matches int

	for k1, _ := range m {
		for _, k2 := range keys {
			if k1 == k2 {
				matches += 1
			}
		}
	}

	return matches == len(keys)
}

// responseProcessor is run internally by Connect() function.
// Tt's responsible for monitoring command responses and notifications
func (b *Bulb) responseProcessor() {
	var buff = make([]byte, 512)
	var resp map[string]interface{}

	for {
		n, err := b.conn.Read(buff)
		if err != nil {
			break
		}

		responses := bytes.Split(buff[:n], []byte{CR, LF})

		for _, r := range responses[:len(responses)-1] {
			resp = make(map[string]interface{})

			err = json.Unmarshal(r, &resp)
			if err != nil {
				log.Printf("OKResponse err: %s\n", r)
				continue
			}

			switch {
			case keysExists(resp, "id", "result"): // Command success
				var unmarshaled OKResponse
				err = json.Unmarshal(r, &unmarshaled)
				if err != nil {
					log.Printf("second unmarshal error: %s\n", r)
				}
				b.results[unmarshaled.id()] <- &unmarshaled
			case keysExists(resp, "id", "error"): // Command failed
				var unmarshaled ERRResponse
				err = json.Unmarshal(r, &unmarshaled)
				if err != nil {
					log.Printf("second unmarshal error: %s\n", r)
				}
				b.results[unmarshaled.id()] <- &unmarshaled
			case keysExists(resp, "method", "params"): // Notification
				// log.Printf("state change%s\n", r)
			default:
				log.Printf("unhandled response: %s\n", r)
			}
		}
	}
	log.Printf("response processor exited\n")
}

// findFirstFreeIntKey finds available (unique) id which will be used as command identifier
func (b *Bulb) findFirstFreeIntKey() (int, error) {
	for i := 0; i < 100; i++ {
		_, ok := b.results[i]
		if !ok {
			return i, nil
		}
	}

	return 0, errors.New("not available")
}
