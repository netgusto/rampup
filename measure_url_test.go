package main

import (
	"errors"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var statusOK = http.StatusOK

func Test_measureURL_normal_case_no_retries(t *testing.T) {
	g := NewGomegaWithT(t)

	statusGetterMock := &URLStatusGetterMock{
		Responses: []URLStatusGetterMockResponse{
			{StatusCode: &statusOK, Error: nil},
		},
	}

	m := measureURL("https://algolia.com", 2, statusGetterMock)
	g.Expect(m.err).ToNot(HaveOccurred())
	g.Expect(m.tries).To(Equal(1))
	g.Expect(m.statusCode).ToNot(BeNil())
	g.Expect(*m.statusCode).To(Equal(http.StatusOK))
}

func Test_measureURL_error_code(t *testing.T) {
	g := NewGomegaWithT(t)

	statusGetterMock := &URLStatusGetterMock{
		Responses: []URLStatusGetterMockResponse{
			{StatusCode: nil, Error: errors.New("TEST ERROR: Something's broken!")},
			{StatusCode: nil, Error: errors.New("TEST ERROR: Something's broken!")},
		},
	}

	m := measureURL("https://whatever", 2, statusGetterMock)
	g.Expect(m.err).To(HaveOccurred())
	g.Expect(m.err.Error()).To(ContainSubstring("TEST ERROR: Something's broken!"))
	g.Expect(m.statusCode).To(BeNil())
	g.Expect(m.tries).To(Equal(2))
	g.Expect(m.duration).ToNot(BeZero())
	g.Expect(m.duration < time.Second).To(BeTrue())
}

func Test_measureURL_works_on_last_try(t *testing.T) {

	syntheticError := errors.New("TRY AGAIN!")
	responses := []URLStatusGetterMockResponse{
		{StatusCode: nil, Error: syntheticError},
		{StatusCode: nil, Error: syntheticError},
		{StatusCode: &statusOK, Error: nil},
	}

	g := NewGomegaWithT(t)
	{
		statusGetterMock := &URLStatusGetterMock{
			Responses: responses,
		}

		m := measureURL("https://whatever", 3, statusGetterMock)
		g.Expect(m.err).ToNot(HaveOccurred())
		g.Expect(m.statusCode).ToNot(BeNil())
		g.Expect(*m.statusCode).To(Equal(http.StatusOK))
		g.Expect(m.tries).To(Equal(3))
	}

	{
		statusGetterMock := &URLStatusGetterMock{
			Responses: responses,
		}

		m := measureURL("https://whatever", 2, statusGetterMock)
		g.Expect(m.err).To(HaveOccurred())
		g.Expect(m.tries).To(Equal(2))
	}
}
