package main

import (
	"miamollie/codeowner-slackbot/gql"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// main is called when a new lambda starts, so don't
// put anything here that needs to happen per request
func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if err := validateSlackRequest(r); err != nil {
			panic("This request wasn't from slack!")
		}
		token := os.Getenv("GITHUB_GQL_AUTH_TOKEN")
		api := os.Getenv("GITHUB_GQL_API")
		c := gql.NewClientWithAuth(api, token)
		Handler(rw, r, c)
	})

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}

func validateSlackRequest(r *http.Request) error {
	// Not going to do this until done testing
	// could do a log.Warn and then pass it through
	// TODO https://api.slack.com/authentication/verifying-requests-from-slack#sdk_support
	return nil
}
