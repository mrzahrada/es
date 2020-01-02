package es

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type PublisherOption func() error

// Streamer -
type Streamer interface {
	Commit(partitionKey string, event Event) error
	Push(ctx context.Context) error
	Clear()
}

// NewStreamer -
// Options:
//	|> stream-name - required
//	|> kinesis-region - optional
//	|> retry policy - optional
func NewStreamer(streamName string) (Streamer, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		return nil, err
	}
	stream := kinesis.New(sess)

	return &kinesisStreamer{
		stream:     stream,
		streamName: streamName,
		commits:    map[string][]Event{},
	}, nil
}

// kinesisStreamer -
type kinesisStreamer struct {
	stream     *kinesis.Kinesis
	streamName string
	commits    map[string][]Event
}

// Commit - implements Streamer.Commit
func (streamer *kinesisStreamer) Commit(partitionKey string, event Event) error {
	if streamer.commits[partitionKey] == nil {
		streamer.commits[partitionKey] = []Event{}
	}
	streamer.commits[partitionKey] = append(streamer.commits[partitionKey], event)
	return nil
}

// Push - implements Streamer.Push
// TODO: I call delete operation, I don't need to call clear anymore
func (streamer *kinesisStreamer) Push(ctx context.Context) error {

	for pk, events := range streamer.commits {
		entries := []*kinesis.PutRecordsRequestEntry{}
		for _, event := range events {
			event.InitializeEvent("user-1", "org-1")
			record, err := marshalEvent(event)
			if err != nil {
				return err
			}
			entries = append(entries, &kinesis.PutRecordsRequestEntry{
				Data:         []byte(record.Data),
				PartitionKey: aws.String(pk),
			})
		}
		delete(streamer.commits, pk)
		if _, err := streamer.stream.PutRecordsWithContext(ctx, &kinesis.PutRecordsInput{
			Records:    entries,
			StreamName: aws.String(streamer.streamName),
		}); err != nil {
			return err
		}
	}

	return nil
}

// Clear stored commits
func (streamer *kinesisStreamer) Clear() {
	streamer.commits = map[string][]Event{}
}

func marshalEvent(v Event) (Record, error) {
	eventType, _ := eventType(v)

	data, err := json.Marshal(v)
	if err != nil {
		return Record{}, ErrMarshalEvent
	}

	entry := jsonEvent{
		Data: data,
		Type: eventType,
	}

	data, err = json.Marshal(entry)
	if err != nil {
		return Record{}, ErrMarshalEvent
	}
	fmt.Println("event:", string(data))

	return Record{
		Data: string(data),
		Type: eventType,
		At:   time.Now().UTC(),
	}, nil
}

func eventType(event interface{}) (string, reflect.Type) {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name(), t
}
