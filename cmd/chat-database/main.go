package main

import (
	"github.com/adrianbrad/chat-v2/pkg/chatdatabase/cmd"
)

func main() {
	command := cmd.NewChatDatabaseCommand()
	command.Execute()
}
