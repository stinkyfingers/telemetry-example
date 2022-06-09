package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"context"
	"os"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

)

const port = "5551"

const appname = "otel_collector"

func main() {
	// set up tracer - export to file
	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	tp := initTracer(f)
	defer tp.Shutdown(context.Background())


	// run server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), middleware(mux)))
}

func middleware (handler http.Handler) http.Handler {
	httpOptions := []otelhttp.Option{
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	}
	return otelhttp.NewHandler(handler, "server1", httpOptions...)
}

func handler1(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(appname).Start(r.Context(), "server1.handler1")
	defer span.End()
	
	fmt.Printf("server1-handler1 traceID %d spanID %d\n", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	
	data, err := getDataFromServer2(ctx)
	if err != nil {
	   	w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	
	w.Write(data)
}

func getDataFromServer2(ctx context.Context) ([]byte, error) {
	ctx, span := otel.Tracer(appname).Start(ctx, "server1.handler1")
	defer span.End()
	
	fmt.Printf("server1-getDataFromServer2 traceID %d spanID %d\n", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	
	resp, err := otelhttp.Get(ctx, "http://localhost:5552/")
	if err != nil {
		return nil, err
	}
	

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)	
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