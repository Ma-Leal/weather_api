package entity

type CEP struct {
	Number string
}

func NewCEP(cep string) (*CEP, error) {

	c := &CEP{Number: cep}
	if err := c.validateCEP(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CEP) validateCEP() error {
	if len(c.Number) != 8 {
		return ErrInvalidCEP
	}
	return nil
}
