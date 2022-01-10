package gql

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type GQLClient interface {
	MakeRequest(vars map[string]interface{}) (gqlResp, error)
}

type client struct {
	Client *graphql.Client
}

func NewClientWithAuth(api string, token string) GQLClient {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token, TokenType: "Bearer"},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return &client{
		Client: graphql.NewClient(api, httpClient),
	}
}

// TODO how to make this nicer and less coupled?
type gqlResp struct {
	Repository struct {
		Description graphql.String
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// Query must be a pointer to a struct that that corresponds to the GraphQL schema
func (c *client) MakeRequest(vars map[string]interface{}) (gqlResp, error) {
	var q gqlResp
	gqlErr := c.Client.Query(context.Background(), &q, vars)
	if gqlErr != nil {
		log.Printf("GQL returned an error: %+v", gqlErr)
		return q, gqlErr
	}

	log.Printf("Response: %+v", q)
	return q, nil
}
