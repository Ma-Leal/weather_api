package usecase

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeatherByCity_Execute(t *testing.T) {
	ctx := context.Background()
	city := "São Paulo"
	key := "k"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"current":{"temp_c":24.1}}`)
	}))
	defer server.Close()

	uc := NewGetWeatherByCity(server.URL, key)
	weather, err := uc.Execute(ctx, city) // agora passando o contexto
	assert.NoError(t, err)
	assert.Equal(t, 24.1, weather.Celsius)
}
