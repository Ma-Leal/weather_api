package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ma-Leal/weather/internal/entity"
	"github.com/Ma-Leal/weather/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WeatherHandler struct {
	GetWeatherByCEP *usecase.GetWeatherByCEP
	OTELTracer      trace.Tracer
}

func NewWeatherHandler(uc *usecase.GetWeatherByCEP, OTELTracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{
		GetWeatherByCEP: uc,
		OTELTracer:      OTELTracer,
	}
}
func (h *WeatherHandler) GetWeatherByCEPHandler(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	otelCtx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	otelCtx, span := h.OTELTracer.Start(otelCtx, "weather-api-request")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "cep is required", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("cep", cep))

	weather, err := h.GetWeatherByCEP.Execute(otelCtx, cep)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get weather by cep")

		switch {
		case errors.Is(err, entity.ErrInvalidCEP):
			http.Error(w, `{"error": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		case errors.Is(err, entity.ErrCEPNotFound):
			http.Error(w, `{"error": "can not find zipcode"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}
