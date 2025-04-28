package usecase

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/Ma-Leal/weather/internal/entity"
)

type GetWeatherByCEP struct {
	getAddress *GetAddressByCEP
	getWeather *GetWeatherByCity
	tracer     trace.Tracer
}

func NewGetWeatherByCEP(getAddress *GetAddressByCEP, getWeather *GetWeatherByCity) *GetWeatherByCEP {
	return &GetWeatherByCEP{
		getAddress: getAddress,
		getWeather: getWeather,
		tracer:     otel.Tracer("GetWeatherByCEP"),
	}
}

func (uc *GetWeatherByCEP) Execute(ctx context.Context, cep string) (*entity.Weather, error) {
	ctx, span := uc.tracer.Start(ctx, "GetWeatherByCEP.Execute")
	defer span.End()

	// Adiciona atributos relevantes ao span
	span.SetAttributes(
		attribute.String("cep", cep),
		attribute.String("component", "usecase"),
	)

	// Executa a busca de endereço
	address, err := uc.getAddress.Execute(ctx, cep)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get address")
		return nil, err
	}

	span.SetAttributes(
		attribute.String("address.city", address.City),
	)

	// Executa a busca de clima
	weather, err := uc.getWeather.Execute(ctx, address.City)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get weather")
		return nil, err
	}

	span.SetAttributes(
		attribute.Float64("weather.temp_c", weather.Celsius),
	)

	return weather, nil
}