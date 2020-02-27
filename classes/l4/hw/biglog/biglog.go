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

func (s searchResult) ToStringArray() (logs []string, err error) {
	for _, v := range s {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		logs = append(logs, string(b))
	}
	return
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

func (l *logger) All(f func(Record) bool) (r searchResult, err error) {
	r = make(searchResult)
	for id, v := range l.data {
		if f(Record(v)) {
			r[id] = Record(v)
		}
	}
	return
}

const logTimeTemplate = "2006-01-02"

func parseLogTime(t string) (time.Time, error) {
	// TODO 时间的解析仍然不够灵活
	return time.ParseInLocation(logTimeTemplate, t, time.Local)
}

type QueryFunc func(Record) (key string, result bool)

//type queryFunc func(Record) bool
func queryFalse(records Record) (key string, result bool) {
	return "", false
}

func (q QueryFunc) Key(key string) QueryFunc {
	q = func(records Record) (string, bool) {
		return key, true
	}
	return q
}

func (q QueryFunc) Or(queries ...QueryFunc) QueryFunc {
	nq := func(records Record) (key string, result bool) {
		for _, q := range queries {
			key, result = q(records)
			if result == true {
				return
			}
		}
		return
	}
	return nq
}

func (q QueryFunc) And(queries ...QueryFunc) QueryFunc {
	nq := func(records Record) (key string, result bool) {
		for _, q := range queries {
			key, result = q(records)
			if result == false {
				return
			}
		}
		return
	}
	return nq
}

func (q QueryFunc) Between(after string, before string) QueryFunc {
	if q == nil {
		return queryFalse
	}
	nFunc := func(records Record) (key string, result bool) {
		pKey, pResult := q(records)
		v, ok := records[pKey]
		if !ok {
			return pKey, false
		}

		switch v.(type) {
		case string:
			// TODO 暂且先把 string 都解析为时间
			vTime, err := parseLogTime(v.(string))

			if err != nil {
				return pKey, false
			}

			afterTime, err := parseLogTime(after)
			if err != nil {
				return pKey, false
			}
			beforeTime, err := parseLogTime(before)
			if err != nil {
				return pKey, false
			}
			afterTime = afterTime.Add(-1)
			beforeTime = beforeTime.Add(1)

			result = vTime.After(afterTime) && vTime.Before(beforeTime)
		case int:
			//TODO
		}

		return pKey, pResult && result
	}
	return nFunc
}

func (q QueryFunc) Equal(value int) QueryFunc {
	if q == nil {
		return queryFalse
	}

	nq := func(records Record) (key string, result bool) {
		pKey, pResult := q(records)
		v, ok := records[pKey]
		if !ok {
			return pKey, false
		}
		switch v.(type) {
		case string:
			// TODO

		case float64:
			// 数字好像都会解析成 float64 ?
			numValue, ok := v.(float64)
			if !ok {
				return pKey, false
			}
			intValue := int(numValue)
			result = intValue == value
		}
		return pKey, pResult && result
	}

	return nq

}

func (q QueryFunc) Exist() QueryFunc {
	if q == nil {
		return queryFalse
	}

	nq := func(records Record) (key string, result bool) {
		pKey, pResult := q(records)
		_, result = records[pKey]
		return pKey, pResult && result
	}

	return nq
}

func (q QueryFunc) End() func(Record) bool {
	f := func(r Record) bool {
		_, result := q(r)
		return result
	}
	return f
}
