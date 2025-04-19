package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Ma-Leal/weather/internal/dto"
	"github.com/Ma-Leal/weather/internal/entity"
)

type GetWeatherByCity struct {
	Url string
	Key string
}

func NewGetWeatherByCity(u string, k string) *GetWeatherByCity {
	return &GetWeatherByCity{Url: u, Key: k}
}

func (w *GetWeatherByCity) Execute(city string) (*entity.Weather, error) {

	params := url.Values{}
	params.Add("key", w.Key)
	params.Add("q", city)
	params.Add("aqi", "no")

	fullUrl := fmt.Sprintf("%s?%s", w.Url, params.Encode())

	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var apiResponse dto.WeatherApiResponse

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return nil, err
	}

	weather := entity.NewWeather(apiResponse.Current.Celsius)
	return weather, nil

}
