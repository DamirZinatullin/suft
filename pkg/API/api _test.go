package api

import (
	//"errors"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

package main

import (
"errors"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
"testing"
)

type mockedHttpClient struct {
	mock.Mock
}

func (m *mockedHttpClient) Get(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func TestAuthenticateGetSuccess(t *testing.T) {
	//mockedClient := new(mockedHttpClient)
	//
	//// Мокаем успешный вызов
	//mockedClient.
	//	On("Get", "https://www.cbr.ru/scripts/XML_daily.asp").
	//	Return("text", nil)
	//
	//rates := CbrRates{httpClient: mockedClient}
	//result, err := rates.Get()
	//assert.NoError(t, err)
	//assert.Equal(t, "text", result)
}

func TestAuthenticate_GetFail(t *testing.T) {
	//mockedClient := new(mockedHttpClient)
	//
	//// Мокаем вызов с ошибкой
	//mockedClient.
	//	On("Get", "https://www.cbr.ru/scripts/XML_daily.asp").
	//	Return("", errors.New("some error"))
	//
	//rates := CbrRates{httpClient: mockedClient}
	//result, err := rates.Get()
	//assert.Error(t, err)
	//assert.Equal(t, "", result)
}

