package usecase

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
)

func Index(ctx context.Context, input string) (map[string]string, error) {
	// meter := global

	ctx, span := otel.Tracer("product-service").Start(ctx, "usecase.index")
	defer span.End()

	_, spanProcessImei := otel.Tracer("product-service").Start(ctx, "usecase.index.processingImei")
	time.Sleep(1 * time.Second)
	spanProcessImei.End()
	span.AddEvent("processingImei(done)")

	_, spanGetCentroid := otel.Tracer("product-service").Start(ctx, "usecase.index.getCentroid")
	time.Sleep(2 * time.Second)
	spanGetCentroid.End()
	span.AddEvent("getCentroid(done)")

	return map[string]string{
		"code":     "200",
		"status":   "OK",
		"endpoint": "index",
	}, nil
}
