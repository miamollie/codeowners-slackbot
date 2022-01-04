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
	log.Printf("REQUEST URL: %#v", r.URL)
	// Validate request: query params, auth header for slack, github token (?)

	// make GQL request (will need to have some kinda private org repo auth check)
	gqlresp := getCodeOwners("99designs", "aws-vault")

	// parse GQL response - and format

	// write response

	io.WriteString(w, gqlresp)
}

// should receive a repo name and org, and return codeowners file content or nil (or empty string, or the actual message?), or an error
func getCodeOwners(org string, repo string) string {
	graphqlClient := graphql.NewClient("https://api.github.com/graphql") //TODO get from .env

	// TODO change to
	gqlQuery := `
    query ($org: String!, $repo: String!)  {
        repository(name: $repo", owner: $org) {
			location1: object(expression: "master:readme.md") {
				... on Blob {
					text
				}
			}
			location2: object(expression: "master:README.md") {
			... on Blob {
					text
				}
			}
		}
    }
`
	req := graphql.NewRequest(gqlQuery)

	req.Var("org", org)
	req.Var("repo", repo)

	req.Header.Set("Cache-Control", "no-cache")

	// run it and capture the response
	var resp RepoResponse
	if err := graphqlClient.Run(context.Background(), req, &resp); err != nil {
		//TODO log the error
		log.Printf("GQL error: %#v", err)

		return "Error making request"
	}

	// if resp.errors != nil {
	// 	return resp.errors.message, nil
	// }
	log.Printf("GQL resp: %#v", resp)
	// TODO split the making of this request and the formatting of the response
	// TODO resp.errors.message might not be null
	return "Made it here"

}

type RepoContents struct {
	text string
}

type RepoResponse struct {
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
