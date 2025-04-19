package handler

import (
	"encoding/json"
	"net/http"

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

	Weather, err := h.GetWeatherByCEP.Execute(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Weather)
}
