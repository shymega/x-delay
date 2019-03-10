package main

import (
	"bytes"
	"github.com/phalaaxx/milter"
	"log"
	"net"
	"time"
)

type XDelayMilter struct {
	milter.Milter
	body *bytes.Buffer
}

func (e *XDelayMilter) Connect(_, _ string, _ uint16, addr net.IP, m *milter.Modifier) (milter.Response, error) {
	log.Printf("Connection received from %s\n", addr.String())

	return milter.RespContinue, nil
}

func (e *XDelayMilter) Helo(name string, m *milter.Modifier) (milter.Response, error) {
	log.Println("Received HELO command.")
	log.Printf("Remote host identified themselves with %s\n", name)

	return milter.RespContinue, nil
}

func (e *XDelayMilter) RcptTo(to string, m *milter.Modifier) (milter.Response, error) {
	log.Printf("Email recipient: %s\n", to)
	return milter.RespContinue, nil
}

func (e *XDelayMilter) BodyChunk(chunk []byte, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *XDelayMilter) Header(k, v string, m *milter.Modifier) (milter.Response, error) {
	if k == "X-Delay" {
		log.Println("X-Delay header detected..")
		log.Printf("Header value: %s\n", v)

		t, _ := time.Parse(time.RFC3339, v)
		result := t.Equal(time.Now())

		if result != true {
			log.Println("Its not time yet to release this message, quarantining..")
			m.Quarantine("Waiting for the right time..")
		}
		return milter.RespAccept, nil
	}
	return milter.RespContinue, nil
}

func (x *XDelayMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	log.Printf("Email from: %s\n", from)

	return milter.RespContinue, nil
}

func (x *XDelayMilter) Body(m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func RunServer(socket net.Listener) {
	init := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
		return &XDelayMilter{},
			milter.OptQuarantine,
			milter.OptNoEOH

	}

	if err := milter.RunServer(socket, init); err != nil {
		panic(err)
	}
}

func main() {
	socket, err := net.Listen("tcp", "127.0.0.1:8896")
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	go RunServer(socket)

	select {}
}
