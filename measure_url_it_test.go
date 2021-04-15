package main

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_IT_measureURL(t *testing.T) {
	g := NewGomegaWithT(t)

	m := measureURL("https://algolia.com", 2, URLStatusGetterReal{})
	g.Expect(m.err).ToNot(HaveOccurred())
	g.Expect(m.tries).To(Equal(1))
	g.Expect(m.statusCode).ToNot(BeNil())
	g.Expect(*m.statusCode).To(Equal(http.StatusOK))
}
