package main

import "fmt"

type x struct {}

func (x x) method() {}

type binary func(int, int) int

func (b binary) addOne(x int) int {
	return b(x, 1)
}

func (b binary) trackTime(x, y int) int {
	t = time.Now()
	defer func(t time.Time) {
		fmt.Println(t.Now() - t)
	}(t)
	return b(x, y)
}


func main() {
	add := binary(func(x, y int) int{
		return x+y
	})

	// 函数调用 f(arg1, arg2...)
	// 类型强转 t(obj)
	add.trackTime()
	fmt.Println(b.addOne(10))
}

type queryFunc func(Record) (bool, bool)

// Open / Close
func (q queryFunc) filter(Record) bool {
	ret, _ := q()
	return ret
}

func (q queryFunc) After(s string) queryFunc {
	return func(Record) bool {
		ret, or = q()
		if or {
			return ret || ..., false
		}
		return ret && ..., false
	}
}

func (q queryFunc) Or() queryFunc {
	return func(Record) bool {
		ret, or = q()
		return ret, true
	}
}

func (q queryFunc) Before(s string) queryFunc {
	return func(Record) bool {
		q()
		....
	}
}
