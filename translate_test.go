package main

import (
	"context"
	"os"
	"testing"
)

var apiKey = os.Getenv("CLOUD_TRANSLATION_API_KEY")

func TestTranslate(t *testing.T) {
	if apiKey == "" {
		t.Fatal("Please set environment variable 'CLOUD_TRANSLATION_API_KEY'.")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, apiKey)
	if err != nil {
		t.Error("failed to create new Translator")
	}

	text, err := client.AutoTranslate(ctx, "こんにちは")
	if text != "Hello" {
		t.Errorf("expected %s, but got %s", "Hello", text)
	}

	text2, err := client.AutoTranslate(ctx, "Hello")
	if text2 != "こんにちは" {
		t.Errorf("expected %s, but got %s", "こんにちは", text)
	}

}
