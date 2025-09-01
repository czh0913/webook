package main

import "fmt"

func NewUser() {
	u := User{}
	fmt.Printf("%v\n", u)
	fmt.Printf("%+v\n", &u)

	up := &User{}
	fmt.Printf("%+v\n", &up)
	up2 := new(User)
	fmt.Printf("%+v\n", &up2)

	u4 := User{Name: "Tom", Age: 0}

	fmt.Printf("u4 : %+v\n", &u4)

}

type User struct {
	Name      string
	FirstName string
	Age       int
}

func (u User) ChangeName(name string) {
	u.Name = name
}

func (u *User) ChangeAge(age int) {
	u.Age = age
}

func ChangeUser() {
	u1 := User{Name: "Tom", Age: 0}
	//u2 := User{Name: "ffff", Age: 12}
	u1.ChangeName("dsfdsdddddd")
	u1.ChangeAge(12121)
	fmt.Println(u1)
}
