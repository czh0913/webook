package main

func Functional4() {
	println("hello , functional4")
}

func Functional5() {
	//新定义了一个方法赋值了fn
	fn := func() string {
		return "hello"
	}

	fn()
}

func UseFunctional4() {
	myFunc := Functional4
	myFunc()
}

func functional8() {
	fn := func() string {
		return "hello this is functional8"
	}()

	println(fn)
}
