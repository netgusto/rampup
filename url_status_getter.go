package main

import "net/http"

type URLStatusGetter interface {
	GetStatus(url string) (*int, error)
}

type URLStatusGetterReal struct{}

func (sg URLStatusGetterReal) GetStatus(url string) (*int, error) {
	var resp *http.Response
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return &resp.StatusCode, nil
}

type URLStatusGetterMockResponse struct {
	StatusCode *int
	Error      error
}

type URLStatusGetterMock struct {
	CalledTimes int
	Responses   []URLStatusGetterMockResponse
}

func (m *URLStatusGetterMock) GetStatus(url string) (*int, error) {
	m.CalledTimes++
	return m.Responses[m.CalledTimes-1].StatusCode, m.Responses[m.CalledTimes-1].Error
}
