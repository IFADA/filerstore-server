package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := os.Create("/sss/1.text")
	if err != nil {
		fmt.Printf("err %s:", err)
	}
}
