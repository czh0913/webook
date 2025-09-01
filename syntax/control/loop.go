package main

import (
	"fmt"
)

func ForLoop() {
	for i := 0; i < 10; i++ {
		println(i)
	}
	for i := 0; i < 10; {
		i++
	}
}

func Loop2() {
	i := 0
	for i < 10 {
		i++
		println(i)
	}

}

func ForArr() {
	arr := [3]int{5, 9, 2}
	for index := range arr {
		println("下标", index, "值", arr[index])
	}
}

func ForMap() {
	m := map[string]int{
		"key1": 1000,
		"key2": 2000,
	}
	// 对map的遍历是随机的
	for k := range m {
		println(k, m[k])
	}
}

func LoopBug() {
	users := []User{
		{
			name: "Tom",
		},
		{
			name: "Jerry",
		},
	}

	m := make(map[string]*User, 2)
	for _, u := range users {
		m[u.name] = &u
	}

	for k, v := range m {
		fmt.Printf("name : %s , user : %v \n", k, v)
	}
}

type User struct {
	name string
}
