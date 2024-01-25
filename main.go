package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	instrument "go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

const metricsInterval = 5 * time.Second

const (
	meterProviderName = "macias_asserts_test"
	metricName        = "events_count"
	svc1Name          = "svc1"
	svc2Name          = "svc2"
	svc3Name          = "svc3"
)

func main() {

	ctx := context.TODO()
	exporter, err := otlpmetrichttp.New(ctx)
	panicOnErr(err)

	svc1, attrs1 := counter(exporter, svc1Name, svc2Name)
	svc2, attrs2 := counter(exporter, svc2Name, svc3Name)
	svc3, attrs3 := counter(exporter, svc3Name, svc1Name)

	for {
		log.Println("increasing counters")
		svc1.Add(ctx, 1, attrs1)
		svc2.Add(ctx, 2, attrs2)
		svc3.Add(ctx, 3, attrs3)
		time.Sleep(5*time.Second)
	}
}

func counter(exporter metric.Exporter, src, dst string) (instrument.Int64Counter, instrument.MeasurementOption) {
	meter := metric.NewMeterProvider(
		metric.WithResource(resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(src),
			semconv.TelemetrySDKLanguageKey.String("go"),
		)),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			metric.WithInterval(metricsInterval))),
	)

	counter, err := meter.Meter(meterProviderName).Int64Counter(metricName)
	panicOnErr(err)

	return counter, instrument.WithAttributes(
		attribute.String("asserts_env", "dev"),
		attribute.String("asserts_site", "dev"),
		attribute.String("namespace", "test"),
		attribute.String("source", src),
		attribute.String("destination", dst),
	)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
