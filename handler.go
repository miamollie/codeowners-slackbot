package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/machinebox/graphql"
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

	io.WriteString(w, gqlresp+" please fricking work")
}

// should receive a repo name and org, and return codeowners file content or nil (or empty string, or the actual message?), or an error
func getCodeOwners(org string, repo string) string {
	graphqlClient := graphql.NewClient("https://api.github.com/graphql") //TODO get from .env. Maybe do the whole gql client setup as part of a service and pass in the endpoint?

	gqlQuery := `
    query ($repo: String!, $org: String!) {
        repository(name: $repo, owner: $org) {
    		location1: object(expression: "master:CODEOWNERS") {
			... on Blob {
				text
			}
    		location2: object(expression: "master:.gitignore/CODEOWNERS") {
			... on Blob {
				text
			}
		}
	}
}
`
	req := graphql.NewRequest(gqlQuery)

	req.Header.Set("Cache-Control", "no-cache")

	// this is the only bit that needs to happen per request
	req.Var("org", org)
	req.Var("repo", repo)


	// run it and capture the response
	var resp RepoResponse
	if err := graphqlClient.Run(context.Background(), req, &resp); err != nil {
		log.Printf("GQL error: %#v", err)

		return "Error making request, soz"
	}

	log.Printf("GQL resp: %#v", resp)
	// TODO split the making of this request and the formatting of the response
	// TODO resp.errors.message might not be null
	return resp.data.repository.location1.text

}

type repoContents struct {
	text string
}

type repoResponse struct {
	data struct {
		repository struct {
			location1 RepoContents
			location2 RepoContents
		}
	}
	errors struct {
		message string
	}
}
