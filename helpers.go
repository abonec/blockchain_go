package main

import "fmt"

func warning(err error, msg string){
	if err != nil {
		fmt.Println(msg + ":", err)
	}
}
