package main

type Inner struct {
}

func (i Inner) Dosomething() {
	println("Inner")
}
func (i Inner) Name() string {
	return "Inner"
}

type Outer struct {
	Inner
}

func (o Outer) Name() string {
	return "Outer"
}

func (o Outer) Dosomething() {
	println("Outer")
}

type OuterPtr struct {
	*Inner
}

type OOOOuter struct {
	Outer
}

func UseInner() {
	var o Outer
	o.Dosomething()

	var op OuterPtr
	op.Dosomething()

	var oopp *OuterPtr
	oopp.Dosomething()
}

func (o Outer) SayHello() {
	println("hello " + o.Name())
}
