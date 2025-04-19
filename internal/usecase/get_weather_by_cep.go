package usecase

import "github.com/Ma-Leal/weather/internal/entity"

type GetWeatherByCEP struct {
	GetAddressByCEP  *GetAddressByCEP
	GetWeatherByCity *GetWeatherByCity
}

func NewGetWeatherByCEP(addresBycep *GetAddressByCEP, weaderByCity *GetWeatherByCity) *GetWeatherByCEP {
	return &GetWeatherByCEP{GetAddressByCEP: addresBycep, GetWeatherByCity: weaderByCity}
}

func (w *GetWeatherByCEP) Execute(cep string) (*entity.Weather, error) {
	address, err := w.GetAddressByCEP.Execute(cep)
	if err != nil {
		return nil, err
	}

	weather, err := w.GetWeatherByCity.Execute(address.City)
	if err != nil {
		return nil, err
	}
	return weather, nil
}
