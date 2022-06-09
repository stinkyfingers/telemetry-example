package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const port = "5552"

func main() {
	tracer.Start()
	defer tracer.Stop()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func handler1(w http.ResponseWriter, r *http.Request) {
	sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
	if err != nil {
	    panic(err)
	}

	span := tracer.StartSpan("server2.handler1", tracer.ChildOf(sctx))
	defer span.Finish()

	fmt.Printf("server2-handler1 traceID %d spanID %d\n", span.Context().TraceID(), span.Context().SpanID())

	w.Write([]byte("server2-handler1-data"))
}