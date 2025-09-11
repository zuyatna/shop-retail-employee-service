package idgen

import "github.com/google/uuid"

type UUIDv7Generator struct{}

func NewUUIDv7Generator() *UUIDv7Generator {
	return &UUIDv7Generator{}
}

func (g *UUIDv7Generator) NewID() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
