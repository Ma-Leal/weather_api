package dto

type WeatherApiResponse struct {
	Current struct {
		Celsius float64 `json:"temp_c"`
	} `json:"current"`
}
