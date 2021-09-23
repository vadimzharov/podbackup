package main

import (
	"os"

	podbackup "github.com/vadimzharov/podbackup/cmd"
)

func main() {
	podbackup.Main(os.Args)
}
