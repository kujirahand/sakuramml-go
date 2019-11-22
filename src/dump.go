package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("dump filename")
		return
	}
	data, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic("Failed to load")
	}
	println("size=", len(data))
	s := ""
	c := ""
	for i := 0; i < len(data); i++ {
		v := data[i]
		s += fmt.Sprintf("%02x ", v)
		c += string(v)
		if i % 8 == 3 {
			s += " "
			c += " "
		}
		if i % 8 == 7 {
			s += "| " + c + "\n"
			c = ""
		}
	}
	println(s + "| " + c)
}


