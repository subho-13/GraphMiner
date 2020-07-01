package main

import "fmt"

func check(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
