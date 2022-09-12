package processor

import (
	"github.com/DusanKasan/parsemail"

	"sync"
)

// derived from
// https://flaviocopes.com/golang-data-structure-queue/

type Email parsemail.Email
type EmailQueue struct {
	lock   sync.Mutex
	emails []Email
}

var emailQueue EmailQueue

func (q *EmailQueue) NewQueue() *EmailQueue {
	q.emails = []Email{}
	return q
}

func (q *EmailQueue) Enqueue(e Email) {
	q.lock.Lock()
	q.emails = append(q.emails, e)
	q.lock.Unlock()
}

func (q *EmailQueue) Dequeue() *Email {
	q.lock.Lock()
	email := q.emails[0]
	q.emails = q.emails[1:len(q.emails)]
	q.lock.Unlock()
	return &email
}

// Do not implement 'Front'
// https://flaviocopes.com/golang-data-structure-queue/

func (q *EmailQueue) IsEmpty() bool {
	return len(q.emails) == 0
}

func (q *EmailQueue) Len() int {
	return len(q.emails)
}
