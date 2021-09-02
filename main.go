package main

import "fmt"

const max = 3

func main() {
	number := [max]int{5, 6, 7}
	var ptrs [max]*int //指针数组
	//将number数组的值的地址赋给ptrs
	for i := 0; i < max; i++ {
		ptrs[i] = &number[i]
	}
	for i, x := range ptrs {
		fmt.Printf("指针数组：索引:%d 值:%d 值的内存地址:%d\n", i, *x, x)
	}
}
