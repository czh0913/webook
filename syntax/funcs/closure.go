package main

func Closure(name string) func() string {
	return func() string {
		return "hello , " + name
	}
}

func Closure1() func() int {
	age := 0

	return func() int {
		age++ //这个age在外部用完之前会一直++，直到重新定义或者用完

		return age
	}
}

func Closure2() func() int {
	k := 1
	return func() int {
		k++
		return k
	}
}
