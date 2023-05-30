package main

import (
	"fmt"
	"os"

	"github.com/provokateurin/rwth-calendar/cmd"
)

func main() {
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
