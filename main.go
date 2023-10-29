package main

import (
	"gocker/command"
	"os"
)

func main() {
	switch os.Args[1] {
	case "run":
		err := command.Run()
		if err != nil {
			panic(err)
		}

	case "init":
		err := command.Init(os.Args[2], os.Args[3], os.Args[3:])
		if err != nil {
			panic(err)
		}
	}
}