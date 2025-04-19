package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCEP_ValidCEP(t *testing.T) {

	cep := "12345678"
	c, err := NewCEP(cep)
	assert.NoError(t, err)
	assert.Equal(t, cep, c.Number)

}

func TestCEP_InvalidCEP(t *testing.T) {
	cep := "12345"
	c, err := NewCEP(cep)
	assert.Error(t, err)
	assert.Empty(t, c)
}
