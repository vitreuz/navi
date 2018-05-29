package postfacto

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPostfacto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postfacto Suite")
}
