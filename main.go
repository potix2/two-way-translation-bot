package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

func translateText(text string) (string, error) {
	apiKey := os.Getenv("CLOUD_TRANSLATE_API_KEY")

	ctx := context.Background()
	client, err := NewClient(ctx, apiKey)
	if err != nil {
		return "", err
	}
	defer client.Close()

	return client.AutoTranslate(ctx, text)
}

// [END translate_translate_text]

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	token := os.Getenv("SLACK_API_TOKEN")
	verificationToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
	event, e := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{verificationToken}))
	if e != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, e
	}
	fmt.Printf("parsed value: %v\n", event)

	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(request.Body), &r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, nil
		}
		return events.APIGatewayProxyResponse{Body: r.Challenge, StatusCode: 200}, nil
	}

	if event.Type == slackevents.CallbackEvent {
		api := slack.New(token)
		postParams := slack.PostMessageParameters{}
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if ev.SubType != "bot_message" {
				translatedText, e := translateText(ev.Text)
				if e != nil {
					return events.APIGatewayProxyResponse{StatusCode: 500}, e
				}
				api.PostMessage(ev.Channel, translatedText, postParams)
			}
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
