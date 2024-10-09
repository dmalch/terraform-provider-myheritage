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

func TestGetIndividualBiography(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	individualId := "individual-760079151-1500318"
	profile, err := GetIndividualBiography(testApiKey, individualId)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.Notes.Data[0].Id).To(BeEquivalentTo("note-760079151-1-1000785"))
}
