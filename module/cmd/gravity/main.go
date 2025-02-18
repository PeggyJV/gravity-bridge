package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/peggyjv/gravity-bridge/module/v6/cmd/gravity/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
