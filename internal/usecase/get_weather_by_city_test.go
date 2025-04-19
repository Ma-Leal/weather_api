package usecase

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeatherByCity_execute(t *testing.T) {
	city := "São Paulo"
	key := "k"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"current":{"temp_c":24.1}}`)
	}))

	defer server.Close()

	uc := NewGetWeatherByCity(server.URL+"?", key)
	wather, err := uc.Execute(city)
	assert.NoError(t, err)
	assert.Equal(t, wather.Celsius, 24.1)

	fmt.Println(uc)
}
