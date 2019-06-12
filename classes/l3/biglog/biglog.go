package biglog

import (
	"crypto/md5"
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
	logMap map[[md5.Size]byte]map[string]interface{}
}
type LogServerErr error

func StartServer() *logServer {
	return &logServer{
		logMap: make(map[[md5.Size]byte]map[string]interface{}),
	}
}

func (l *logServer) Close() error {
	return nil
}

func (l *logServer) Log(r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	info := make(map[string]interface{})
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return err
	}
	//info["CreatedAt"] = time.Now()
	// TODO 感觉用 map 不太合适
	// 取 md5 来作为 key 可以避免重复
	key := md5.Sum(bytes)
	l.logMap[key] = info
	return nil
}
func convertMapToJsonString(jsonMap map[[md5.Size]byte]map[string]interface{}) ([]string, error) {
	results := make([]string, 0)
	for _, v := range jsonMap {
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		results = append(results, string(bytes))
	}
	return results, nil
}

func parseLogTime(t string) (time.Time, error) {
	// TODO 时间的解析仍然不够灵活
	return time.ParseInLocation(logTimeTemplate, t, time.Local)
}

func filterTimeCondition(jsonMap map[[md5.Size]byte]map[string]interface{}, key string,
	conditionType string, conditionValue time.Time) (map[[md5.Size]byte]map[string]interface{}, error) {
	var f func(time.Time, time.Time) bool
	switch conditionType {
	case ">":
		f = func(t1 time.Time, t2 time.Time) bool {
			t2 = t1.Add(-1)
			return t1.After(t2)
		}
	case "<":
		f = func(t1 time.Time, t2 time.Time) bool {
			t2 = t2.Add(1)
			return t1.Before(t2)
		}
	case "=":
		f = func(t1 time.Time, t2 time.Time) bool {
			return t1.Equal(t2)
		}
	default:
		return nil, fmt.Errorf(" %v 不是正确的操作符", conditionType)
	}
	result := make(map[[md5.Size]byte]map[string]interface{})
	for logK, logV := range jsonMap {
		v, ok := logV[key]
		if !ok {
			continue
		}
		logTimeValue, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("对应时间字段 %v 不是 string 类型", key)
		}
		logTime, err := parseLogTime(logTimeValue)
		if err != nil {
			return nil, fmt.Errorf("对应字段 %v 在转换为 time.Time 时出错 %v", key, err.Error())
		}
		if f(logTime, conditionValue) {
			result[logK] = logV
		}
	}
	return result, nil
}

// TODO 暂时只支持 int
func filterNumCondition(jsonMap map[[md5.Size]byte]map[string]interface{}, key string,
	conditionType string, conditionValue int) (map[[md5.Size]byte]map[string]interface{}, error) {
	var f func(float64, float64) bool
	switch conditionType {
	case ">":
		f = func(v1 float64, v2 float64) bool {
			return v1 > v2
		}
	case "<":
		f = func(v1 float64, v2 float64) bool {
			return v1 < v2
		}
	case "=":
		f = func(v1 float64, v2 float64) bool {
			return v1 == v2
		}
	default:
		return nil, fmt.Errorf(" %v 不是正确的操作符", conditionType)
	}
	result := make(map[[md5.Size]byte]map[string]interface{})
	for logK, logV := range jsonMap {
		v, ok := logV[key]
		if !ok {
			continue
		}
		// 这里 v 是 float64 类型。。。暂时还不知道为什么
		intValue, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("对应字段 %v 不是数字类型", key)
		}
		if f(intValue, float64(conditionValue)) {
			result[logK] = logV
		}
	}
	return result, nil
}

// TODO 实现的不好
func (l *logServer) Search(conditions map[string][]interface{}) ([]string, error) {
	result := l.logMap
	for key, conditions := range conditions {

		if len(conditions)&1 == 1 {
			return nil, errors.New("条件和参数不匹配")
		}
		var conditionType string
		var conditionValue interface{}

		for i := 0; i < len(conditions); i++ {
			// 偶数位为条件类型 目前有 > < =
			if i&1 == 0 {
				conditionType = conditions[i].(string)
				continue
			}
			conditionValue = conditions[i]
			switch conditionValue.(type) {
			case string:
				// TODO 暂且先把 string 都解析为时间
				t, err := parseLogTime(conditionValue.(string))
				if err != nil {
					return nil, err
				}
				result, err = filterTimeCondition(result, key, conditionType, t)
				if err != nil {
					return nil, err
				}
			case int:
				// 这里写成 := 会导致重新创建一个 result
				var err error
				result, err = filterNumCondition(result, key, conditionType, conditionValue.(int))
				if err != nil {
					return nil, err
				}
			}

		}
	}

	return convertMapToJsonString(result)
}
