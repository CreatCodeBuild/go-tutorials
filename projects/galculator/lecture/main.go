package main // 主、主要的

import "fmt" // format 格式

// 主函数、程序从这里开始
func main() {
	// result 结果
	加, 减, 乘, 除, 余数, 商 := compute简单版(-100, 200)
	fmt.Println("compute简单版的结果：", 加, 减, 乘, 除, 余数, 商)

	// Print 打印
	// ln line 行

	fmt.Println("compute升级版", compute升级版(1, 2, "+"))
	fmt.Println("compute升级版", compute升级版(1, 2, "-"))
	fmt.Println("compute升级版", compute升级版(1, 2, "*"))
	fmt.Println("compute升级版", compute升级版(1, 2, "/"))
	//                  0  1  2  3  4   5     6
	floats := []float64{1, 2, 3, 4, 10, -3.2, 3.1415926}
	fmt.Println("compute升级版2", floats)
	fmt.Println("compute升级版2", floats[0])
	fmt.Println("compute升级版2", floats[6])
	fmt.Println("compute升级版2", compute升级版2(floats))
}

// 函数
// 函数名 compute
// func -- function 函数
// func 函数名(零个或者多个参数) 可选的返回值 {
//     函数具体的运算内容
// }
// input 输入
// func compute(input string) int {

// }

func compute简单版(
	input1 int, input2 int,
) (int, int, int, int, int, float64) {
	return input1 + input2, input1 - input2, input1 * input2, input1 / input2, input1 % input2, float64(input1) / float64(input2)
}

func compute升级版(input1, input2 float64, 运算 string) float64 {
	// if 如果
	if 运算 == "+" {
		return input1 + input2
	} else if 运算 == "-" {
		return input1 - input2
	} else if 运算 == "*" {
		return input1 * input2
	} else if 运算 == "/" {
		return input1 / input2
	} else {
		return -1
	}
}

// 接片 Slice
func compute升级版2(input []float64) float64 {
	// 为了、只要
	// for 起始条件；继续条件；迭代更新 {

	// }

	// len -- length 长度
	// 变量名 := 初始值
	// var 变量名 类型 = 值
	var sum float64
	for i := 0; i < len(input); i += 1 {
		fmt.Println(i, sum)
		sum += input[i]
		fmt.Println(i, sum)
	}
	return sum
}

func 指数(底数 float64，幂 int) {
	
}
