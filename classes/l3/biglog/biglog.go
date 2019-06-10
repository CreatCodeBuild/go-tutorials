package biglog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

type LogInfo struct {
	Number    int       `json:"number"`
	Time      LogTime   `json:"time"`
	CreatedAt time.Time `json:"-"`
}

type LogTime time.Time

const logTimeJsonTemplate = "\"2006-01-02\""
const logTimeTemplate = "2006-01-02"

func (t LogTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("%s", time.Time(t).Format(logTimeJsonTemplate))
	return []byte(stamp), nil
}

func (t *LogTime) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	lt, err := time.ParseInLocation(logTimeJsonTemplate, s, time.Local)
	if err != nil {
		return err
	}
	*t = LogTime(lt)
	return nil
}

type logServer struct {
	logMap map[int]LogInfo
}
type LogServerErr error

func StartServer() *logServer {
	return &logServer{logMap: map[int]LogInfo{}}
}

func (l *logServer) Close() error {
	return nil
}

func (l *logServer) Log(r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	info := LogInfo{}
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return err
	}
	info.CreatedAt = time.Now()
	// TODO 感觉用 map 不太合适
	l.logMap[info.Number] = info
	return nil
}

// TODO 实现的不好
func (l *logServer) Search(conditions map[string]string) ([]string, error) {
	logM := l.logMap
	if value, ok := conditions["after"]; ok {
		after, err := time.ParseInLocation(logTimeTemplate, value, time.Local)
		if err != nil {
			return nil, errors.New("after time error")
		}
		after = after.Add(-1)
		tmpMap := map[int]LogInfo{}
		for k, v := range logM {
			t := time.Time(v.Time)
			if t.After(after) {
				tmpMap[k] = v
			}
		}
		logM = tmpMap
	}

	if value, ok := conditions["before"]; ok {
		before, err := time.ParseInLocation(logTimeTemplate, value, time.Local)
		if err != nil {
			return nil, errors.New("before time error")
		}
		before = before.Add(1)
		tmpMap := map[int]LogInfo{}
		for k, v := range logM {
			t := time.Time(v.Time)
			if t.Before(before) {
				tmpMap[k] = v
			}
		}
		logM = tmpMap
	}
	results := make([]string, 0)
	for _, v := range logM {
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		results = append(results, string(bytes))
	}
	return results, nil
}
