package biglog

import (
	"encoding/json"
	"io"
	"math/rand"
	"time"
)

type record map[string]interface{}

type logger struct {
	data   map[int64]record
	nextID func() int64
}

// Record is used for user to define their own filter function.
type Record map[string]interface{}

type searchResult map[int64]Record

func (s searchResult) String() (string, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func StartServer() logger {
	rand.Seed(time.Now().UnixNano())
	return logger{
		data: make(map[int64]record),
		nextID: func() int64 {
			return time.Now().Unix() + rand.Int63()
		},
	}
}

func (l *logger) Log(data io.Reader) error {
	id := l.nextID()
	r := make(record)
	err := json.NewDecoder(data).Decode(&r)
	if err != nil {
		return err // todo: handle it
	}
	l.data[id] = r
	return nil
}

func (l *logger) All(f Query) (r searchResult, err error) {
	r = make(searchResult)
	for id, v := range l.data {
		if f(Record(v)) {
			r[id] = Record(v)
		}
	}
	return
}

type Query func(Record) bool

func (q Query) Key(key string) Query {

	return q
}

func (q Query) Or(queries ...Query) Query {

	return q
}

func (q Query) And(queries ...Query) Query {

	return q
}

func (q Query) Between(after string, before string) Query {

	return q
}

func (q Query) Equal(value int) Query {

	return q
}

func (q Query) Exit() Query {

	return q
}
