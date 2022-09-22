package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"time"
)

// Event is the event structure for eventbridge events
// see https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-events.html
type Event struct {
	Version    string                 `json:"version"`
	Id         string                 `json:"id"`
	DetailType string                 `json:"detail-type"`
	Source     string                 `json:"source"`
	Account    string                 `json:"account"`
	Time       time.Time              `json:"time"`
	Region     string                 `json:"region"`
	Resources  []string               `json:"resources"`
	Detail     map[string]interface{} `json:"detail"`
}

// EventBus is the event bus structure
type EventBus struct {
	Name       string
	DetailType string
	Source     string
	Data       any
}

// Fire an event to the eventbridge
func (e *EventBus) Fire(sess *session.Session) error {
	if e.Name == "" {
		return fmt.Errorf("event bus name is required")
	}

	jsonData, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	if e.DetailType == "" {
		e.DetailType = "default"
	}

	if e.Source == "" {
		e.Source = "default"
	}

	evb := eventbridge.New(sess)
	entries := []*eventbridge.PutEventsRequestEntry{{
		Detail:       aws.String(string(jsonData)),
		DetailType:   aws.String(e.DetailType),
		EventBusName: aws.String(e.Name),
		Source:       aws.String(e.Source),
	}}

	_, err = evb.PutEventsWithContext(context.Background(), &eventbridge.PutEventsInput{Entries: entries})

	return err
}
