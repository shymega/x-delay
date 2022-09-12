package processor

import (
	"github.com/DusanKasan/parsemail"

	"sync"
)

type EmailQueue struct {
	lock   sync.Mutex
	emails []parsemail.Email
}

var emailQueue EmailQueue
