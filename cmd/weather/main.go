package main

import (
	"net/http"

	"github.com/Ma-Leal/weather/configs"
	"github.com/Ma-Leal/weather/internal/handler"
	"github.com/Ma-Leal/weather/internal/usecase"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	ucGetAddress := usecase.NewGetAddressByCEP(cfg.CepApiUrl)
	ucGetWeather := usecase.NewGetWeatherByCity(cfg.WeatherApiUrl, cfg.WeatherApiKey)
	ucGetWeatherByCEP := usecase.NewGetWeatherByCEP(ucGetAddress, ucGetWeather)

	weatherHandler  := handler.NewWeatherHandler(ucGetWeatherByCEP)

	http.HandleFunc("/", weatherHandler.GetWeatherByCEPHandler)
	http.ListenAndServe(":8080", nil)
}
