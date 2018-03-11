package main

import (
	_ "github.com/lib/pq"
	"github.com/m1kola/shipsterbot/internal/cli"
)

func main() {
	cli.Execute()
}
