package examples

import (
	"fmt"

	"github.com/jwillp/gogram/telegram"
)

const (
	appID    = 6
	appHash  = "YOUR_APP_HASH"
	botToken = "YOUR_BOT_TOKEN"
)

func main() {
	// Create a new client
	client, _ := telegram.NewClient(telegram.ClientConfig{
		AppID:    appID,
		AppHash:  appHash,
		LogLevel: telegram.LogInfo,
	})

	// Connect to the server
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// Authenticate the client using the bot token
	if err := client.LoginBot(botToken); err != nil {
		panic(err)
	}

	// new conversation
	conv, _ := client.NewConversation("username or id", 30) // 30 is the timeout in seconds
	defer conv.Close()
	_, err := conv.SendMessage("Hello, Please reply to this message")
	if err != nil {
		panic(err)
	}
	resp, err := conv.GetResponse() // wait for the response
	// resp, err := conv.GetReply() // wait for the reply
	// conv.MarkRead() // mark the conversation as read
	// conv.WaitEvent() // wait for any custom update
	if err != nil {
		panic(err)
	}
	// Print the response
	fmt.Println("Response:", resp.Text())
}
