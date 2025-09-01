package main

func main() {

	//// _ 在这里是不接受返回值
	//
	//fn := Closure1()
	//println(fn())
	//println(fn())
	//println(fn())
	//println(fn())
	//println(fn())
	//println(fn())
	//fn = Closure1()
	//println(fn())
	//println(fn())
	//println(fn())
	//println(fn())
	//println(fn())
	f := Closure2()
	println(f())
	println(f())
	d := Closure2()
	println(d())
	println(d())
	println(d())
}
