package main

import (
	_ "github.com/lib/pq"
	"github.com/m1kola/shipsterbot/cmd"
)

func main() {
	cmd.Execute()
}
