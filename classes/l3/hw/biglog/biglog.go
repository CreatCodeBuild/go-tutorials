package biglog

import (
	"time"
	"encoding/json"
	"io"
	"fmt"
	"math/rand"
)
// StartServer returns a started logger.
func StartServer() logger {
	rand.Seed(time.Now().UnixNano())
	return logger{
		data: make(map[int64]record),
		nextID: func() int64 {
			return time.Now().Unix() + rand.Int63()
		},
	}
}

type record map[string]interface{}

type logger struct {
	data map[int64]record
	nextID func() int64
}

func (l *logger) Log(data io.Reader) error {	
	id := l.nextID()
	r := make(record)
	err := json.NewDecoder(data).Decode(&r)
	if err != nil {
		return err // todo: handle it
	}
	l.data[id] =  r
	return nil
}

// Record is used for user to define their own filter function.
type Record map[string]interface{}

type searchResult map[int64]Record

func (s searchResult) String() (string, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (l *logger) All(f func(Record) bool) (r searchResult, err error) {
	r = make(searchResult)
	for id, v := range l.data {
		if f(Record(v)) {
			r[id] = Record(v)
		}
	}
	return
}


