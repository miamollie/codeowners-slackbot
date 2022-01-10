package main

import (
	"fmt"
	"io"
	"log"
	"miamollie/codeowner-slackbot/gql"
	"net/http"
	"strings"

	"github.com/shurcooL/graphql"
)

// handler is the function called when the lambda is invoked
// handler function signature (context, TIn where TIn can be unmarshalled by Go)
// return value must be either error or TOut where TOut can marshalled

func Handler(w http.ResponseWriter, r *http.Request, c gql.GQLClient) {
	// Validate request: query params, auth header for slack, github token (..?)
	log.Printf("Query Params: %#v", r.URL.RawQuery) //needs to be decoded and then split to get the value
	// TODO use CHI instead so you can get at the query params more easily
	p := strings.Split(r.URL.RawQuery, "/")
	if len(p) != 2 {
		log.Printf("Didn't get enough query params")
	}

	log.Printf("split query params are %+v ", p)
	// make GQL request (will need to have some kinda private org repo auth check)
	gqlresp := getCodeOwners(c, "miamollie", "codeowners-slackbot")

	io.WriteString(w, string(gqlresp)+" Description!")
}



// should receive a repo name and org, and return codeowners file content or nil (or empty string, or the actual message?), or an error
func getCodeOwners(c gql.GQLClient, owner string, repo string) graphql.String {
	if len(owner) == 0 || len(repo) == 0 {
		return graphql.String(fmt.Sprintf("Invalid arguments: %s, %s", owner, repo))
	}

	vars := map[string]interface{}{
		"owner": graphql.String(owner),
		"name":  graphql.String(repo),
	}

	resp, e := c.MakeRequest(vars)

	if e != nil {
		log.Printf("error: %v", e)
		//TODO handle other error cases
		return "Something went wrong! TODO, handle specific error cases"
	}

	if resp.Repository.Description != "" {
		return resp.Repository.Description
	}

	// Error cases
	// - auth " Sorry, you are not authorised to read files from this repository"
	// - file doesn't exist "No Codeowners file could be found in this repository"
	// something else "Unknown error :("
	// SUccess: contents of the file at one of the locations

	return "Sorry, could not find a file at that location"
}
