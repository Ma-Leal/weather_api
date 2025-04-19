package dto

type ViaCepResponse struct {
	Cep  string `json:"cep"`
	City string `json:"localidade"`
}
