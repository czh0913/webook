package main

// T 是类型参数，名字叫T ，约束是 any ，等于没有约束
type List[T any] interface {
	Add(idx int, t T)
	Append(t T)
}

func UseList() {
	var l List[int]
	l.Append(1)
}

type LinkedList[T any] struct {
	head *node[T]
	t    T
	e    **T
}
type node[T any] struct {
	val T
}
