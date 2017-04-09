package bolt

import (
	"encoding/json"
	"github.com/inkah-trace/server"
	pb "github.com/inkah-trace/daemon/protobuf"
	"fmt"
)

var _ server.TraceService = &TraceService{}

type TraceService struct {
	client *Client
}

func (ts *TraceService) CreateEvent(event *pb.ForwardedEvent) error {
	tx, err := ts.client.db.Begin(true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tx.Rollback()

	tb := tx.Bucket([]byte("Traces"))
	tesb := tx.Bucket([]byte("TraceEvents"))

	teb := tesb.Bucket([]byte(event.TraceId))
	if teb == nil {
		teb, err = tesb.CreateBucket([]byte(event.TraceId))
		if err != nil {
			fmt.Println(err)
			return err
		}
		ti := server.TraceInfo{
			TraceId: event.TraceId,
			Hostname: event.Hostname,
			Timestamp: event.Timestamp,
		}
		tij, err := json.Marshal(&ti)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tb.Put(i64tob(event.Timestamp), tij)
	}


	id, err := teb.NextSequence()
	if err != nil {
		fmt.Println(err)
		return err
	}

	ej, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		return err
	}

	teb.Put(itob(int(id)), ej)

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (ts *TraceService) TraceList() ([]*server.TraceInfo, error) {
	tx, err := ts.client.db.Begin(true)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	// Assume bucket exists and has keys
	b := tx.Bucket([]byte("Traces"))

	c := b.Cursor()

	traces := make([]*server.TraceInfo, 0)
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var ti server.TraceInfo
		json.Unmarshal(v, &ti)
		traces = append(traces, &ti)
	}

	return traces, nil
}

func (ts *TraceService) Trace(id string) ([]*server.EventInfo, error) {
	tx, err := ts.client.db.Begin(true)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	// Assume bucket exists and has keys
	tb := tx.Bucket([]byte("TraceEvents"))
	b := tb.Bucket([]byte(id))

	c := b.Cursor()

	events := make([]*server.EventInfo, 0)
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var fe pb.ForwardedEvent
		json.Unmarshal(v, &fe)

		e := server.EventInfo{
			SpanId: fe.SpanId,
			ParentSpanId: fe.ParentSpanId,
			Timestamp: fe.Timestamp,
			EventType: fe.EventType,
		}

		events = append(events, &e)
	}

	return events, nil
}