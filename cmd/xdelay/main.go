package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	log "github.com/inconshreveable/log15"
)

var (
	bindAddr string
	domain   string
	mtaAddr  string
	mtaPort  string
)

type Backend struct{}
type User struct{}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.User, error) {
	return &User{}, nil
}

func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.User, error) {
	return &User{}, nil
}

func (u *User) Send(from string, to []string, r io.Reader) error {
	log.Info("Received message from MTA, now processing.")

	var b []byte
	var err error

	if b, err = ioutil.ReadAll(r); err != nil {
		return err
	} else {
		byteIO := bytes.NewReader(b)

		email, err := parsemail.Parse(byteIO)
		if err != nil {
			log.Error(err.Error())
		}

		if len(email.Header.Get("X-Delay-TS")) > 0 {
			log.Info("X-Delay-TS header detected.")
			log.Info("Queuing for further processing.")

			err = ioutil.WriteFile("/home/dzr/message.dat", b, 0644)

			return nil
		}

		err = smtp.SendMail("localhost:10026", nil, from, to, byteIO)
		if err != nil {
			log.Error("Failed to send email back to MTA.",
				log.Ctx{"error": err.Error()})
		}
	}

	return nil
}

func (u *User) Logout() error {
	return nil
}

func init() {
	flag.StringVar(&bindAddr, "bindAddr",
		"127.0.0.1:10025",
		"Address to listen on for SMTP requests.")

	flag.StringVar(&domain, "domain",
		"localhost",
		"Domain to use for SMTP server.")

	flag.StringVar(&mtaAddr, "mtaAddr",
		"localhost",
		"MTA address for message injection.")

	flag.StringVar(&mtaPort, "mtaPort",
		"10026",
		"MTA port for message injection.")

	flag.Parse()
}

func main() {
	bkd := &Backend{}
	s := smtp.NewServer(bkd)

	s.Addr = bindAddr
	s.Domain = domain
	s.MaxIdleSeconds = 150
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 20
	s.Strict = true
	s.AuthDisabled = true

	log.Info("Starting server now.",
		log.Ctx{"bindAddress": s.Addr})
	if err := s.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
