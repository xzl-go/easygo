// Package tracing 提供了基于 OpenTelemetry 的分布式追踪功能
// 支持跨服务调用链路的追踪和监控
package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Tracer 是追踪器
// 负责创建和管理追踪跨度（Span）
type Tracer struct {
	tracer     trace.TracerProvider // 追踪器提供者
	tracerName string               // 服务名称
}

// NewTracer 创建一个新的追踪器
// serviceName: 服务名称，用于标识追踪来源
func NewTracer(serviceName string) *Tracer {
	// 创建标准输出导出器，用于调试
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	// 创建追踪器提供者
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // 使用批处理器导出追踪数据
	)

	// 设置全局追踪器提供者
	otel.SetTracerProvider(tp)

	return &Tracer{
		tracer:     tp,
		tracerName: serviceName,
	}
}

// StartSpan 开始一个新的追踪跨度
// ctx: 上下文
// spanName: 跨度名称
// 返回新的上下文和追踪跨度
func (t *Tracer) StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := t.tracer.Tracer(t.tracerName)
	return tracer.Start(ctx, spanName)
}

// EndSpan 结束追踪跨度
// span: 要结束的追踪跨度
func (t *Tracer) EndSpan(span trace.Span) {
	span.End()
}

// Shutdown 关闭追踪器
// ctx: 上下文，用于控制关闭超时
// 返回关闭错误（如果有）
func (t *Tracer) Shutdown(ctx context.Context) error {
	if tp, ok := t.tracer.(*sdktrace.TracerProvider); ok {
		return tp.Shutdown(ctx)
	}
	return nil
}
