package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Ma-Leal/weather/internal/entity"
	"github.com/Ma-Leal/weather/internal/usecase"
)

type WeatherHandler struct {
	GetWeatherByCEP *usecase.GetWeatherByCEP
}

func NewWeatherHandler(uc *usecase.GetWeatherByCEP) *WeatherHandler {
	return &WeatherHandler{GetWeatherByCEP: uc}
}

func (h *WeatherHandler) GetWeatherByCEPHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "cep is required", http.StatusBadRequest)
		return
	}

	weather, err := h.GetWeatherByCEP.Execute(cep)
	fmt.Println(err)

	if err != nil {
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
