package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

// handler is the function called when the lambda is invoked
// handler function signature (context, TIn where TIn can be unmarshalled by Go)
// return value must be either error or TOut where TOut can marshalled

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Query Params: %#v", r.URL.RawQuery) //needs to be decoded and then split to get the value
	// Validate request: query params, auth header for slack, github token (..?)

	// make GQL request (will need to have some kinda private org repo auth check)
	gqlresp := getCodeOwners("miamollie", "codeowners-slackbot")

	// parse GQL response - and format

	io.WriteString(w, string(gqlresp)+" Description!")
}

// should receive a repo name and org, and return codeowners file content or nil (or empty string, or the actual message?), or an error
func getCodeOwners(owner string, repo string) graphql.String {
	token := os.Getenv("GITHUB_GQL_AUTH_TOKEN")
	api := os.Getenv("GITHUB_GQL_API")
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token, TokenType: "Bearer"},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	graphqlClient := graphql.NewClient(api, httpClient)

	var query struct {
		Repository struct {
			Description graphql.String
		} `graphql:"repository(name: $name, owner: $owner)"`
	}

	vars := map[string]interface{}{
		"owner": graphql.String(owner),
		"name":  graphql.String(repo),
	}

	gqlErr := graphqlClient.Query(context.Background(), &query, vars)
	if gqlErr != nil {
		log.Printf("GQL returned an error: %+v", gqlErr)
	}

	log.Printf("Response: %+v", query)

	return query.Repository.Description
}
