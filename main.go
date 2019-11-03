package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

type Book struct {
	ID       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name  string
	Books []int
}

type Comment struct {
	Body string
}

func populate() []Book {
	author := &Author{Name: "Robert C. Martin", Books: []int{2}}

	var books []Book
	books = append(books, Book{
		ID:     1,
		Title:  "Clean Code",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "This book is awesome!"},
		},
	})
	books = append(books, Book{
		ID:     2,
		Title:  "Clean Architecture",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "This book is also awesome!"},
		},
	})

	return books
}

func main() {
	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"books": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)
	var bookType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Book",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	books := populate()

	fields := graphql.Fields{
		"book": &graphql.Field{
			Type:        bookType,
			Description: "Get Book By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, book := range books {
						if int(book.ID) == id {
							return book, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(bookType),
			Description: "Get Book List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return books, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})

	err = http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal(err)
	}
}
