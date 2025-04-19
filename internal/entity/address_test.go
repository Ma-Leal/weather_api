package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {

	cep := "12345678"
	c, err := NewCEP(cep)
	assert.NoError(t, err)

	city := "São Paulo"

	add := NewAddress(c, city)

	assert.Equal(t, c.Number, add.Cep.Number)
	assert.Equal(t, city, add.City)

}
