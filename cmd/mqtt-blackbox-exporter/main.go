package main

import (
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
