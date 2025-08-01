package main

import (
	"<<!.ProjectName!>>/internal/server"
	"<<!.ProjectName!>>/internal/utils"
)

func main() {
	srv := server.New()
	
	if err := srv.Start(); err != nil {
		utils.Logger.Error("Server failed to start", "error", err)
		panic(err)
	}
}