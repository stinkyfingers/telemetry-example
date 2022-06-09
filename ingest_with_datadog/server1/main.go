package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const port = "5551"

func main() {
	tracer.Start()
	defer tracer.Stop()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func handler1(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), "server1.handler1")
	defer span.Finish()
	
	fmt.Printf("server1-handler1 traceID %d spanID %d\n", span.Context().TraceID(), span.Context().SpanID())
	
	data, err := getDataFromServer2(ctx)
	if err != nil {
	   	panic(err)
	}
	
	w.Write(data)
}

func getDataFromServer2(ctx context.Context) ([]byte, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "server1.getDataFromServer2")
	defer span.Finish()

	fmt.Printf("server1-getDataFromServer2 traceID %d spanID %d\n", span.Context().TraceID(), span.Context().SpanID())

	req, err := http.NewRequest("GET", "http://localhost:5552/", nil)
	if err != nil {
	    return nil, err
	}
	err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
	if err != nil {
	    return nil, err
	}
	
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
	    return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)	
}