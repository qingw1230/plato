package logger

import (
	"context"
	"testing"
	"time"

	ptrace "github.com/qingw1230/plato/common/prpc/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/qingw1230/plato/common/config"
)

func TestLogger(t *testing.T) {
	config.Init("../../im.yaml")
	NewLogger(WithLogDir("/Users/www/logs"))
	InfoCtx(context.Background(), "info test")
	DebugCtx(context.Background(), "debug test")
	WarnCtx(context.Background(), "warn test")
	ErrorCtx(context.Background(), "error test")
	time.Sleep(1 * time.Second)
}

func TestTraceLog(t *testing.T) {
	config.Init("../../im.yaml")
	NewLogger(WithLogDir("/Users/www/logs"))
	ptrace.StartAgent()
	defer ptrace.StopAgent()

	tr := otel.GetTracerProvider().Tracer(ptrace.TraceName)
	ctx, span := tr.Start(context.Background(), "logger-trace", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	InfoCtx(ctx, "test")
}
