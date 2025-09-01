package main

func YourName(name string, arr ...string) {

}

func CallYourName() {
	YourName("11")
	YourName("11", "22")
	YourName("11", "22", "33")

	arr := []string{"a", "b", "c"}
	YourName("dsf", arr...)
}
