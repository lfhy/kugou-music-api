package kugou_test

import (
	"context"
	"testing"

	kg "github.com/lfhy/kugou-music-api"
)

func TestRootPackageFacade(t *testing.T) {
	client, err := kg.New(kg.WithLite(true), kg.WithCookie(map[string]string{"token": "test-token"}))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if client == nil {
		t.Fatal("New() returned nil client")
	}

	client2, err := kg.NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client2 == nil {
		t.Fatal("NewClient() returned nil client")
	}

	if got, ok := kg.RouteByIdentifier("search"); !ok || got != kg.RouteSearch {
		t.Fatalf("RouteByIdentifier(search) = %q, %v", got, ok)
	}
	if got, ok := client.RouteByIdentifier("search"); !ok || got != kg.RouteSearch {
		t.Fatalf("client.RouteByIdentifier(search) = %q, %v", got, ok)
	}

	if len(kg.Endpoints()) == 0 {
		t.Fatal("Endpoints() returned empty list")
	}
	if len(kg.APIList) == 0 {
		t.Fatal("APIList is empty")
	}

	_ = kg.SearchRequest{Keywords: "test", Page: 1, Pagesize: 10}
	_ = kg.PersonalRadioRequest{Mode: kg.PersonalRadioHeart}
	_ = kg.SongPlayURLRequest{Hash: "hash"}

	var searchFn func(context.Context, kg.SearchRequest) (*kg.SearchResponse, error) = client.Search
	if searchFn == nil {
		t.Fatal("client.Search method is not exposed from root package")
	}
}
