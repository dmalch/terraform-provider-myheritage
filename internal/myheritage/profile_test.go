package myheritage

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetProfile(t *testing.T) {
	RegisterTestingT(t)

	profile, err := GetProfile("", "")

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
}
