package entity

type Weather struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewWeather(c float64) *Weather {
	return &Weather{
		Celsius:    c,
		Fahrenheit: c*1.8 + 32,
		Kelvin:     c + 273,
	}
}
