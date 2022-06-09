package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

const port = "5552"

const appname = "otel_collector"

var propagator propagation.TextMapPropagator

func main() {
	// set up tracer - export to file
	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	tp := initTracer(f)
	defer tp.Shutdown(context.Background())
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), middleware(mux)))
}

func middleware (handler http.Handler) http.Handler {
	return otelhttp.NewHandler(handler, "server2")
}

func handler1(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(appname).Start(r.Context(), "server2.handler1")
	defer span.End()

	fmt.Printf("server2-handler1 traceID %d spanID %d\n", span.SpanContext().TraceID(), span.SpanContext().SpanID())

	w.Write([]byte("server2-handler1-data"))
}

func initTracer(r io.Writer) *trace.TracerProvider {
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(r), stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}