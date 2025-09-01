package main

func Map() {
	m1 := map[string]int{
		"key1": 123,
	}
	m1["hello"] = 232
	//预估容量，不传默认16
	m2 := make(map[string]int, 12)
	m2["key2"] = 456
	val, ok := m1["dam"]
	if ok {
		println(val)
	}

	val = m1["dam"]
	println("对应的值是", val)
}
