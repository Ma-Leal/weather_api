package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/Ma-Leal/weather/internal/dto"
	"github.com/Ma-Leal/weather/internal/entity"
)

type GetWeatherByCity struct {
	Url    string
	Key    string
	client *http.Client
	tracer trace.Tracer
}

func NewGetWeatherByCity(u string, k string) *GetWeatherByCity {
	return &GetWeatherByCity{
		Url: u,
		Key: k,
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		tracer: otel.Tracer("GetWeatherByCity"),
	}
}

func (w *GetWeatherByCity) Execute(ctx context.Context, city string) (*entity.Weather, error) {
	ctx, span := w.tracer.Start(ctx, "GetWeatherByCity.Execute")
	defer span.End()

	params := url.Values{}
	params.Add("key", w.Key)
	params.Add("q", city)
	params.Add("aqi", "no")

	fullUrl := fmt.Sprintf("%s?%s", w.Url, params.Encode())

	span.SetAttributes(
		attribute.String("http.url", fullUrl),
		attribute.String("city", city),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create HTTP request")
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := w.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to make HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
	)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to read response body")
		return nil, err
	}

	var apiResponse dto.WeatherApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unmarshal response")
		return nil, err
	}

	weather := entity.NewWeather(apiResponse.Current.Celsius)

	span.SetAttributes(
		attribute.Float64("weather.temp_c", weather.Celsius),
	)

	return weather, nil
}
