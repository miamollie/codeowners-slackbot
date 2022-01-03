package main

import (
	"fmt"
	"os"
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/machinebox/graphql"
)

func Handler() {

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")

	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(token, slack.OptionDebug(true))
	// Create the Slack attachment that we will send to the channel
	attachment := slack.Attachment{
		Pretext: "Super Bot Message",
		Text:    "some text",
		// Color Styles the Text, making it possible to have like Warnings etc.
		Color: "#36a64f",
		// Fields are Optional extra data!
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}
	// PostMessage will send the message away.
	// First parameter is just the channelID, makes no sense to accept it
	_, timestamp, err := client.PostMessage(
		"", //TODO use the channelID from the request
		// uncomment the item below to add a extra Header to the message, try it out :)
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent at %s", timestamp)
}

// what should this package look like?
// in lambda the gql client could be shared

func Todo() {
	graphqlClient := graphql.NewClient("https://<GRAPHQL_API_HERE>") //TODO from env github api
	req := graphql.NewRequest(`
    query ($key: String!) {
        items (id:$key) {
            field1
            field2
            field3
        }
    }
`)

	// set any variables
	req.Var("key", "value")

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// run it and capture the response
	var resp interface{} //todo be more specific
	if err := graphqlClient.Run(context.Background(), req, &resp); err != nil {
		panic(fmt.Errorf("error: %s", err)) //print to logs and return nil
	}
	fmt.Println(resp) //return resp
}
