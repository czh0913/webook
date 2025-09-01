package main

func Defer() {
	println("1")
	defer func() {
		println("第一个 defer")
	}()
	println("2")
	defer func() {
		println("第二个 defer")
	}()
	println("3")

	//先按顺序输出defer之外内容，倒着执行所有defer
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}() //这里是我没有传入值，最后执行的时候读入了一个i
	i = 1
}

func DeferClosure2() {
	i := 0
	defer func(i int) {
		println(i)
	}(i) //这里是相当于我传入了当前的i但是我程序后面在执行
	i = 1
}
