package myheritage

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetIndividual(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	individualId := "individual-760079151-1500318"
	profile, err := GetIndividual(testApiKey, individualId)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.FamilyGroups[0].Type).To(BeEquivalentTo("parent"))
}
