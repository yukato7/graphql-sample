package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

func resolveID(p graphql.ResolveParams) (interface{}, error) {
	return p.Args["id"], nil
}

func resolveName(p graphql.ResolveParams) (interface{}, error) {
	return "hoge", nil
}

var q graphql.ObjectConfig = graphql.ObjectConfig{
	Name: "query",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: resolveID,
		},
		"name": &graphql.Field{
			Type:    graphql.String,
			Resolve: resolveName,
		},
	},
}

var m graphql.ObjectConfig = graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "Params",
				Fields: graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.Int,
					},
					"game_profile": &graphql.Field{
						Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "state",
							Fields: graphql.Fields{
								"game_id": &graphql.Field{
									Type: graphql.Int,
								},
								"player_name": &graphql.Field{
									Type: graphql.String,
								},
							},
						}),
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Todo add insert process
				return User{
					ID: 1,
					GameProfile: Profile{
						GameID:     1,
						PlayerName: "焔の錬金術師",
					},
				}, nil
			},
		},
	},
}

type User struct {
	ID          uint    `json:"id"`
	GameProfile Profile `json:"game_profile"`
}

type Profile struct {
	GameID     uint   `json:"game_id"`
	PlayerName string `json:"player_name"`
}

var schemaConfig graphql.SchemaConfig = graphql.SchemaConfig{
	Query:    graphql.NewObject(q),
	Mutation: graphql.NewObject(m),
}

var schema, _ = graphql.NewSchema(schemaConfig)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func handler(w http.ResponseWriter, r *http.Request) {
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(r.Body)
	query := bufBody.String()

	result := executeQuery(query, schema)
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
