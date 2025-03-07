package main

import (
	"fmt"
)

func multiply(a, b int) int {
	product := a * b
	return product
}

func add(a, b int) int {
	sum := a + b
	return sum
}

func main() {
	fmt.Println("HELLOW world")
	var intNum int16 = 200
	myNum := -3.14
	a, b := 10, "hello"
	var (
		c int    = 1
		d string = "example"
	)
	if len(d) > len(b) {
		fmt.Println("1st true")
	} else if a > c {
		fmt.Println("2nd true")
	} else {
		fmt.Println("else true")
	}

	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	count := 0
	for count < 5 {
		fmt.Println("Count:", count)
		count++
	}

	result := multiply(4, 6)
	fmt.Println("Multiplied:", result)

	resulta := add(5, 3)
	fmt.Println("Added:", resulta)

	fmt.Println(intNum)
	fmt.Println(myNum)

	person := map[string]string{
		"name": "Alice",
		"city": "New York",
	}
	person["age"] = "30"
	fmt.Println(person)

	numbers := []int{1, 2, 3}
	numbers = append(numbers, 4)
	fmt.Println(numbers)

	array := [5]int{10, 20, 30, 40, 50}
	slice := array[1:4]
	fmt.Println(slice)
}
