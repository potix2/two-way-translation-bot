package main

import (
	"context"
	"html"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type Translator struct {
	client *translate.Client
}

// [START translate_translate_text]

var langMap = map[language.Tag]language.Tag{
	language.Japanese: language.English,
	language.English:  language.Japanese,
}

func NewClient(ctx context.Context, apiKey string) (*Translator, error) {
	//apiKey := os.Getenv("CLOUD_TRANSLATE_API_KEY")
	c, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Translator{
		client: c,
	}, nil
}

func (t *Translator) detectLanguage(ctx context.Context, text string) (language.Tag, error) {
	detections, err := t.client.DetectLanguage(ctx, []string{text})
	if err != nil {
		return language.Und, err
	}

	return (detections[0][0]).Language, nil

}

func (t *Translator) AutoTranslate(ctx context.Context, text string) (string, error) {
	lang, err := t.detectLanguage(ctx, text)
	if err != nil {
		return "", err
	}
	return t.Translate(ctx, langMap[lang], text)
}

func (t *Translator) Translate(ctx context.Context, lang language.Tag, text string) (string, error) {
	resp, err := t.client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", err
	}

	return html.UnescapeString(resp[0].Text), nil
}

// [END translate_translate_text]

func (t *Translator) Close() {
	t.client.Close()
}
