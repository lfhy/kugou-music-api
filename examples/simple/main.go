package main

import (
	"context"
	"fmt"
	"log"

	kg "github.com/lfhy/kugou-music-api"
)

func main() {
	client, err := kg.New()
	if err != nil {
		log.Fatalf("init sdk failed: %v", err)
	}

	resp, err := client.Search(context.Background(), kg.SearchRequest{
		Keywords: "周杰伦",
		Page:     1,
		Pagesize: 10,
	})
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	fmt.Printf("status: %d\n", resp.Status)
	fmt.Printf("body: %s\n", string(resp.RawBody))
}
