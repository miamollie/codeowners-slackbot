package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

// handler is the function called when the lambda is invoked
// handler function signature (context, TIn where TIn can be unmarshalled by Go)
// return value must be either error or TOut where TOut can marshalled

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Query Params: %#v", r.URL.RawQuery) //needs to be decoded and then split to get the value
	// Validate request: query params, auth header for slack, github token (?)

	// make GQL request (will need to have some kinda private org repo auth check)
	gqlresp := getCodeOwners("miamollie", "codeowners-slackbot")

	// parse GQL response - and format

	io.WriteString(w, string(gqlresp)+" Description!")
}

// should receive a repo name and org, and return codeowners file content or nil (or empty string, or the actual message?), or an error
func getCodeOwners(owner string, repo string) graphql.String {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_GQL_AUTH_TOKEN"), TokenType: "Bearer"},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	graphqlClient := graphql.NewClient(os.Getenv("GITHUB_GQL_API"), httpClient)

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
		fmt.Println(gqlErr)
	}

	return query.Repository.Description
}
