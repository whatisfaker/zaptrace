# ZapTrace

使用zap作为基础 logger，符合opentracing规范的微服务tracing用基础库, 参考jagger的例子并复用了部分example代码

使用Jagger作为context tracing的Opentracing实现

## Log

使用zaplog作为基础log然后包装为context相关/无关的log

### 普通log(不追踪context)

```golang
logger := log.NewStdLogger(log.InfoLevel)
logger.Normal().Debug("DebugInfo", zap.String("field", "value"))
logger.Normal().Info("Info", zap.String("field", "value"))
// output is
// {"level":"info","ts":"2019-06-13T13:19:32.318+0800","caller":"log/log_test.go:25","msg":"Info"}
// log level = InfoLevel 所以 debug信息不会输出
```

### 追踪log(追踪context)

```golang
// 根据ctx创建rootSpan 一般不用自己创建，我们的调用都是从http开始，所以都会被封装的middle做掉了
// 这里的ctx是里面不带span的
span, ctx := opentracing.StartSpanFromContext(ctx, "operation_name")
defer span.Finish()

// 纯新创建rootSpan，上下文无关
sp := opentracing.StartSpan("operation_name")
defer sp.Finish()

// 创建子span
// 从上下文中获取Parent Span
var span opentracing.Span
if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
    // 创建子span
    span = opentracing.StartSpan(
        "operation_name",
        opentracing.ChildOf(parentSpan.Context())
    )
// else 这里是做保障，没取到父亲就自己为最root的span(正常情况下, 如果没有取到父亲说明追踪链出问题了)
} else {
    span = opentracing.StartSpan(
        "operation_name"
    )
}
defer span.Finish()

// 设置一些追踪字段 参考本文tracing部分的spec
span.SetTag  (尽量使用spec内的标准推荐tag)
span.SetBaggageItem (慎用)

// 额外增加Trace(ctx) 其他和Normal一致
logger := log.NewStdLogger(log.InfoLevel)
logger.Trace(ctx).Debug("DebugInfo", zap.String("field", "value"))
logger.Trace(ctx).Info("Info", zap.String("field", "value"))
// 数据除了本地输出进标准log外，还会传入opentracing对应的span log中
```

> 注意：对于追踪log, 用于opentracing的logField会默认使用掉"level" 和 "message"所以自己写zapField的时候避免使用这两个关键字

## Tracing

### 使用Jaeger开源库，使用opentracing标准

[OpenTracing语义标准](https://github.com/opentracing-contrib/opentracing-specification-zh/blob/master/specification.md)

主要关注 span tag 和 log field的推荐定义
对Go开发来说，参考[span tag](https://github.com/opentracing/opentracing-go/blob/master/ext/tags.go)

### 使用tracing服务，环境变量设置

主要关注[采样](https://www.jaegertracing.io/docs/1.12/sampling/)

和 [参数定义](https://github.com/jaegertracing/jaeger-client-go#environment-variables)

*JAEGER_SERVICE_NAME*环境变量已经在代码内强制定义了，所以此环境变量无效

### http 服务的追踪

1. [标准库](https://github.com/opentracing-contrib/go-stdlib)
2. [gin的封装](https://github.com/opentracing-contrib/go-gin) 一般使用标准库改装可以应用所有web框架，但是gin本身的responseWriter被改造了，导致标准库的middleware无法追踪responseWriter的statusCode,一直返回0而报错，所以用gin这个重新封装

### gRPC 服务追踪

[gRPC](https://github.com/opentracing-contrib/go-grpc)

### AMQP 消息队列追踪 (参考)

[AMQP](https://github.com/opentracing-contrib/go-amqp) 这里只做参考，具体实现要根据公司内的rabbitMQ进行改造
