package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Ma-Leal/weather/internal/dto"
	"github.com/Ma-Leal/weather/internal/entity"
)

type GetAddressByCEP struct {
	Url string
}

func NewGetAddressByCEP(url string) *GetAddressByCEP {
	return &GetAddressByCEP{Url: url}
}

func (a *GetAddressByCEP) Execute(cep string) (*entity.Address, error) {
	c, err := entity.NewCEP(cep)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/%s/json/", a.Url, c.Number)
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r dto.ViaCepResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.New("unmarshal failed")
	}

	add := entity.NewAddress(
		c,
		r.City,
	)

	return add, nil
}
