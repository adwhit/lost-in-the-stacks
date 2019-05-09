package main

import ( "fmt" )

func a(x int) map[string]int {
	var y = 100;
	return b(x, y)
}

func b(x int, y int) map[string]int {
	m := make(map[string]int);
	m["x"] = x;
	m["y"] = y;
	return m
}

func main() {
	amap := a(10);
	for k, v := range amap {
		fmt.Printf("Key: %s\tValue: %v\n", k, v);
	}
}
