package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lfhy/kugou-music-api/sdk"
)

func main() {
	client, err := sdk.New()
	if err != nil {
		log.Fatalf("init sdk failed: %v", err)
	}

	resp, err := client.Search(context.Background(), sdk.SearchRequest{
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
