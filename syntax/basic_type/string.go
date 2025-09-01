package main

import "unicode/utf8"

func String() {
	println("He said : \"Hello Go ! \" ")
	println("Hellok, Go !")
	println(len("你好"))
	println(utf8.RuneCountInString("你好5165"))

}
