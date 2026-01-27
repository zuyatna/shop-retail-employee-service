package usecase_test

import "github.com/stretchr/testify/mock"

type MockIDGenerator struct {
	mock.Mock
}

func (m *MockIDGenerator) NewID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
