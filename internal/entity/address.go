package entity

type Address struct {
	Cep  *CEP   `json:"cep"`
	City string `json:"cidade"`
}

func NewAddress(cep *CEP, city string) *Address {
	return &Address{Cep: cep, City: city}
}
