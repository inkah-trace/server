package server

import (
	"github.com/graphql-go/graphql"
	"github.com/inkah-trace/collector"
	"github.com/inkah-trace/collector/sqlite"
)

var traceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Traces",
		Description: "",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.String,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if trace, ok := p.Source.(*collector.Trace); ok {
						return trace.Id, nil
					}
					return nil, nil
				},
			},
			"start": &graphql.Field{
				Type:        graphql.Int,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ss := sqlite.SpanService{}
					var traceId string
					if trace, ok := p.Source.(*collector.Trace); ok {
						traceId = trace.Id
					}
					spans := ss.Spans(traceId)

					minTime := int64(99999999999999)
					for _, s := range spans {
						if s.Start.Unix() < minTime {
							minTime = s.Start.Unix()
						}
					}
					return minTime, nil
				},
			},
			"end": &graphql.Field{
				Type:        graphql.Int,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ss := sqlite.SpanService{}
					var traceId string
					if trace, ok := p.Source.(*collector.Trace); ok {
						traceId = trace.Id
					}
					spans := ss.Spans(traceId)

					maxTime := int64(0)
					for _, s := range spans {
						if s.End.Unix() > maxTime {
							maxTime = s.End.Unix()
						}
					}
					return maxTime, nil
				},
			},
			"spans": &graphql.Field{
				Type:        graphql.NewList(spanType),
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ss := sqlite.SpanService{}
					var traceId string
					if trace, ok := p.Source.(*collector.Trace); ok {
						traceId = trace.Id
					}
					spans := ss.Spans(traceId)
					return spans, nil
				},
			},
		},
	},
)

var spanType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Spans",
		Description: "",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.String,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if span, ok := p.Source.(*collector.Span); ok {
						return span.Id, nil
					}
					return nil, nil
				},
			},
			"programName": &graphql.Field{
				Type:        graphql.String,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if span, ok := p.Source.(*collector.Span); ok {
						return span.ProgramName, nil
					}
					return nil, nil
				},
			},
			"hostname": &graphql.Field{
				Type:        graphql.String,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if span, ok := p.Source.(*collector.Span); ok {
						return span.Hostname, nil
					}
					return nil, nil
				},
			},
			//"trace": &graphql.Field{
			//	Type:        traceType,
			//	Description: "",
			//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//		ts := sqlite.TraceService{}
			//		fmt.Println(p)
			//		id, _ := p.Args["TraceId"].(string)
			//		return ts.Trace(id), nil
			//	},
			//},
			//"parentSpan": &graphql.Field{
			//	Type:        spanType,
			//	Description: "",
			//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//		ss := sqlite.SpanService{}
			//		fmt.Println(p)
			//		id, _ := p.Args["ParentId"].(string)
			//		return ss.Span(id), nil
			//	},
			//},
			"start": &graphql.Field{
				Type:        graphql.Int,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if span, ok := p.Source.(*collector.Span); ok {
						return span.Start.Unix(), nil
					}
					return nil, nil
				},
			},
			"end": &graphql.Field{
				Type:        graphql.Int,
				Description: "",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if span, ok := p.Source.(*collector.Span); ok {
						return span.End.Unix(), nil
					}
					return nil, nil
				},
			},
		},
	},
)

var fields = graphql.Fields{
	"traces": &graphql.Field{
		Type: graphql.NewList(traceType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ts := sqlite.TraceService{}
			traces := ts.Traces()
			return traces, nil
		},
	},
}

var rootQuery = graphql.ObjectConfig{Name: "InkahQuery", Fields: fields}

var Schema = graphql.SchemaConfig{
	Query: graphql.NewObject(rootQuery),
}
