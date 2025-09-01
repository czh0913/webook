package main

func IfElse(age int) string {
	if age >= 18 {
		println("成年")

	} else if age < 18 {
		println("未成年")
	} else {
		println("不是人")
	}
	return ""
}

func IfNewVariable(start int, end int) string {
	if distance := end - start; distance > 100 {
		return "too far"
	} else if distance > 60 {
		return "a little"
	} else {
		return "ok"
	}
}
