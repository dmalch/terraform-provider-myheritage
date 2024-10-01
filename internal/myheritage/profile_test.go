package myheritage

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetProfile(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	profile, err := GetProfile("", "")

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
}
