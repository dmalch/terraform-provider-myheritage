package myheritage

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetProfileHeader(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	profileId := "profile-760079151-1500318-0"
	profile, err := GetProfileHeader(testApiKey, profileId)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.FirstName).To(BeEquivalentTo("Иона Герасимович"))
}

func TestGetProfileDetails(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	profileId := "profile-760079151-1500318-0"
	profile, err := GetProfileDetails(testApiKey, profileId)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.FamilyGroups[0].Type).To(BeEquivalentTo("parent"))
}
