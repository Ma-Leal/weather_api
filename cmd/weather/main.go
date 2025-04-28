package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Ma-Leal/weather/configs"
	"github.com/Ma-Leal/weather/internal/handler"
	"github.com/Ma-Leal/weather/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	shutdown, err := initProvider("weather-api", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = shutdown(context.Background())
	}()

	tracer := otel.Tracer("weather-api")

	weatherHandler := initDependencies(cfg, tracer)

	http.HandleFunc("/", weatherHandler.GetWeatherByCEPHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func initDependencies(cfg *configs.Conf, tracer trace.Tracer) *handler.WeatherHandler {
	ucGetAddress := usecase.NewGetAddressByCEP(cfg.CepApiUrl)
	ucGetWeather := usecase.NewGetWeatherByCity(cfg.WeatherApiUrl, cfg.WeatherApiKey)
	ucGetWeatherByCEP := usecase.NewGetWeatherByCEP(ucGetAddress, ucGetWeather)

	weatherHandler := handler.NewWeatherHandler(ucGetWeatherByCEP, tracer)
	return weatherHandler
}
