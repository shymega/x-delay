package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"time"

	"github.com/DusanKasan/parsemail"
	smtp "github.com/emersion/go-smtp"
	log "github.com/inconshreveable/log15"
)

var (
	bindAddr string
	domain   string
	mtaAddr  string
	mtaPort  string
)

type Backend struct{}
type Session struct{}

// Login callback function for go-smtp.
// We use this to accept any SMTP connection. Authentication should be
// handled by the MTA, not us.
func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	log.Debug("Connection accepted as authenticated login.")
	return &Session{}, nil
}

// AnonymousLogin callback function for go-smtp.
// We use this to accept any SMTP connection. Authentication should be
// handled by the MTA, not us.
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	log.Debug("Connection accepted as anonymous login.")
	return &Session{}, nil
}

// Mail (MAIL FROM) callback function for go-smtp.
// We use this to handle MAIL FROM commands from SMTP client(s).
func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Debug("MAIL FROM command accepted.")
	return nil
}

// Rcpt (RCPT TO) callback function for go-smtp.
// We use this to handle RCPT TO commands from SMTP client(s).
func (s *Session) Rcpt(to string) error {
	log.Debug("RCPT TO command accepted.")
	return nil
}

// checkError is a quick 'n dirty function to handle errors.
func checkError(e error) {
	if e != nil {
		log.Crit(e.Error())
	}
}

func (s *Session) Data(r io.Reader) error {
	log.Info("Email received from MTA; now processing.")

	var b []byte
	var err error

	log.Debug("Reading email into byte array.")
	if b, err = ioutil.ReadAll(r); err != nil {
		log.Crit(err.Error()) /* we should handle this gracefully */
	}

	log.Debug("Parsing email into struct.")
	email, err := parsemail.Parse(bytes.NewReader(b))
	if err != nil {
		log.Crit(err.Error()) /* we should handle this gracefully */
	}

	log.Debug("Searching for X-Delay-TS header.")
	if email.Header.Get("X-Delay-TS") != "" {
		log.Info("Found X-Delay-TS header.")
		log.Info("Sending for further processing!")

		err = ioutil.WriteFile("/home/dzr/message.dat", b, 0644)

		return nil
	}

	log.Debug("Injecting email back to MTA.")
	err = smtp.SendMail(
		fmt.Sprintf("%s:%s",
			mtaAddr, mtaPort),
		nil,
		fmt.Sprintf("<%s>",
			email.From[0].Address),
		concatToCcBccHeader(email.To, email.Cc, email.Bcc),
		bytes.NewReader(b))

	if err != nil {
		log.Crit(err.Error()) /* we should handle this gracefully */
	}

	return nil
}

func concatToCcBccHeader(to []*mail.Address, cc []*mail.Address, bcc []*mail.Address) []string {
	var final []string

	log.Debug("Begin processing of To, Cc, and Bcc headers.")

	for _, e := range to {
		final = append(final, e.Address)
	}

	for _, e := range cc {
		final = append(final, e.Address)
	}

	for _, e := range bcc {
		final = append(final, e.Address)
	}

	log.Debug("Processing finalised.")

	return final
}

func (s *Session) Reset() {
	log.Debug("RSET commad accepted.")
}

func (s *Session) Logout() error {
	log.Debug("LOGOUT command accepted.")
	return nil
}

func init() {
	log.Info("X-Delay initialising.")

	flag.StringVar(&bindAddr, "bindAddr",
		"localhost:10025",
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
	log.Debug("Creating new SMTP server.")
	bkd := &Backend{}
	s := smtp.NewServer(bkd)

	log.Debug("Configuring SMTP server parameters..")

	s.Addr = bindAddr
	s.AuthDisabled = true
	s.Domain = domain
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.ReadTimeout = 20 * time.Second
	s.Strict = true
	s.WriteTimeout = 20 * time.Second

	log.Info("Initialisation complete!")

	log.Info("Starting server now.",
		log.Ctx{"bindAddress": s.Addr})
	if err := s.ListenAndServe(); err != nil {
		log.Crit(err.Error()) /* we should handle this gracefully */
	}
}
