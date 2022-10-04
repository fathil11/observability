package main

import (
	"context"
	"fathil/golang-tracing/delivery"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	serviceName        = "product-service"
	serviceEnvironment = "development"
	serviceVersion     = "1.0"
)

func main() {
	err := initTracerProvider("localhost:4318")
	if err != nil {
		log.Fatal(err)
	}

	s := gin.Default()

	s.GET("product", delivery.Index)
	s.GET("product/:id", delivery.Show)

	s.Run(":5000")
}

func initTracerProvider(endpoint string) error {
	ctx := context.Background()

	// Create meter exporter
	meterExp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(endpoint),
	)
	if err != nil {
		return err
	}

	// Create meter provider
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(meterExp, metric.WithInterval(2*time.Second))),
	)

	global.SetMeterProvider(meterProvider)

	// Create trace client
	tracerClient := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endpoint),
	)

	// Create trace exporter
	tracerExp, err := otlptrace.New(ctx, tracerClient)
	if err != nil {
		return err
	}

	// Create trace resource
	tracerRes, err := resource.New(ctx,
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("service.environment", serviceEnvironment),
			attribute.String("service.version", serviceVersion),
		),
	)

	if err != nil {
		return err
	}

	// Create trace span processor
	traceProcessor := sdktrace.NewBatchSpanProcessor(tracerExp)

	// Create trace provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(tracerRes),
		sdktrace.WithSpanProcessor(traceProcessor),
	)

	otel.SetTracerProvider(tracerProvider)

	return nil
}
