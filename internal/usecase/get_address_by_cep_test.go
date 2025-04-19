package usecase

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAddressByCEP_Execute(t *testing.T) {
	cep := "12345678"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"cep": "12345678", "localidade": "São Paulo"}`)
	}))
	defer server.Close()

	uc := NewGetAddressByCEP(server.URL)
	address, err := uc.Execute(cep)
	assert.NoError(t, err)
	assert.Equal(t, cep, address.Cep.Number)
	assert.Equal(t, "São Paulo", address.City)
}
