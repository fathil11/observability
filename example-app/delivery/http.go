package delivery

import (
	"fathil/golang-tracing/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
)

var meter = global.Meter("product-meter")
var requestCount, _ = meter.SyncInt64().Counter(
	"product-service/request_counts",
	instrument.WithDescription("The number of requests received"),
)

func Index(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()

	requestCount.Add(ctx, 1)

	ctx, span := otel.Tracer("product-service1").Start(ctx, "handler.index")
	defer span.End()

	time.Sleep(500 * time.Millisecond)

	result, err := usecase.Index(ctx, "test")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	gCtx.JSON(200, result)
}

func Show(gCtx *gin.Context) {

}
