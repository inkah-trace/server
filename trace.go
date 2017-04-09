package server

import (
	pb "github.com/inkah-trace/daemon/protobuf"
)

type TraceService interface {
	CreateEvent(event *pb.ForwardedEvent) error
	TraceList() ([]*TraceInfo, error)
	Trace(id string) ([]*EventInfo, error)
}

type TraceInfo struct {
	TraceId string `json:"traceId"`
	Hostname string `json:"hostname"`
	Timestamp int64 `json:"timestamp"`
}

type EventInfo struct {
	SpanId string `json:"spanId"`
	ParentSpanId string `json:"parentSpanId"`
	Timestamp int64 `json:"timestamp"`
	EventType pb.EventType `json:"eventType"`
}