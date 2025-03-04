package main

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"strings"
)

type Brand struct {
	Name  string
	Terms string
}

type ESResponse struct {
	Hits struct {
		Hits []struct {
			Source Brand `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

var esClient *elasticsearch.Client

func main() {
	var err error
	esClient, err = elasticsearch.NewClient(elasticsearch.Config{})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	log.Println(esClient.Info())

	initializeIndex()
	insertDocuments()
	initializeServer()
}

func initializeIndex() {
	_, _ = esClient.Indices.Delete([]string{"brands_index"})

	_, err := esClient.Indices.Create("brands_index")
	if err != nil {
		log.Fatalf("Error creating the index: %s", err)
	}

	log.Println("Index created successfully")
}

func insertDocuments() {
	brands := []Brand{
		{
			Name:  "Apple",
			Terms: "Apple, IPOD, Iphone, Macbook",
		},
		{
			Name:  "Samsung",
			Terms: "Samsung, Galaxy, S21, S20",
		},
	}

	for _, brand := range brands {
		insertDocument(brand)
	}
}

func insertDocument(brand Brand) {
	data, _ := json.Marshal(brand)
	_, err := esClient.Index("brands_index", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Error inserting the document: %s", err)
	}
}

func initializeServer() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		terms := r.URL.Query().Get("terms")
		if terms == "" {
			http.Error(w, "Missing 'terms' query parameter", http.StatusBadRequest)
			return
		}

		result, err := findAllDocuments(terms)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})

	http.HandleFunc("/brand", func(w http.ResponseWriter, r *http.Request) {
		var brand Brand
		if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		insertDocument(brand)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Document inserted successfully"))
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func findAllDocuments(terms string) ([]byte, error) {
	query := `{
		"query": {
			"match": {
				"Terms": {
					"query": "` + terms + `",
					"fuzziness": "AUTO"
				}
			}
		},
		"size": 1
	}`

	search, err := esClient.Search(
		esClient.Search.WithIndex("brands_index"),
		esClient.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer search.Body.Close()

	var esResponse ESResponse
	if err := json.NewDecoder(search.Body).Decode(&esResponse); err != nil {
		return nil, err
	}

	if len(esResponse.Hits.Hits) == 0 {
		return nil, nil
	}

	brand := esResponse.Hits.Hits[0].Source
	result, err := json.Marshal(brand)
	if err != nil {
		return nil, err
	}

	return result, nil
}
