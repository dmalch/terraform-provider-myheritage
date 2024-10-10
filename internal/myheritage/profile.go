package myheritage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
)

type ProfileHeaderResponse struct {
	Data struct {
		Profile ProfileHeader `json:"profile"`
	} `json:"data"`
}

type ProfileHeader struct {
	FirstName  string           `json:"first_name"`
	LastName   string           `json:"last_name"`
	Individual IndividualHeader `json:"individual"`
}

type IndividualHeader struct {
	Id        string `json:"id"`
	BirthDate struct {
		Text string `json:"text"`
	} `json:"birth_date"`
	DeathDate struct {
		Text string `json:"text"`
	} `json:"death_date"`
	BirthPlace string `json:"birth_place"`
	DeathPlace string `json:"death_place"`
	Age        struct {
		Text string `json:"text"`
	} `json:"age"`
	Gender       string `json:"gender"`
	CauseOfDeath string `json:"cause_of_death"`
}

func CreateProfile(apiKey, name, description string) (string, error) {
	return "", nil
}

func GetProfileHeader(apiKey, profileId string) (*ProfileHeader, error) {
	// Create a buffer to hold the multipart form-data
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)

	// Add the query part
	rawQuery := `{profile(id:"` + profileId + `",lang:"EN"){name first_name last_name gender age_group age{text}personal_photo{...personal_photo_fragment}is_prefer_user can_current_user_view_discoveries can_current_user_manage_photos can_current_user_edit_personal_photo can_current_user_invite_individual tabs{name total counters}recent_individuals{data{...history_fragment}}favorite_individuals{data{...history_fragment}}individual{...individual_fragment}user{...user_fragment}site_membership{member_id site_id can_user_contact_member member_joined_date member_last_visit_date role_sentence{text}}tree{is_imported_using_family_search_sync}site{name}}}fragment history_fragment on Individual{id name gender age_group lifespan personal_photo{...personal_photo_fragment}tree_relationship{description}link_in_profile_page}fragment individual_fragment on Individual{id name first_name gender personal_photo{...personal_photo_fragment}is_privatized religious_name former_name namesake alternate_names birth_date{text}birth_place death_date{text}is_alive is_likely_deceased death_place burial_place cause_of_death is_cause_of_death_holocaust age{text}tree_relationship{description is_blood_relative is_biological_blood_relative blood_relative_description hour_glass_color_code is_path_to_self is_path_cannot_decide_if_related}can_edit link_in_tree link_in_pedigree_tree link_in_fan_view link_in_research_this_person link_template_in_edit_profile}fragment personal_photo_fragment on Photo{thumbnails(thumbnail_size:"136x136c"){url}}fragment user_fragment on User{name first_name crown_status country country_code birth_date{text}age{text}age_group_in_years show_age is_public is_privatized nickname created_time}`
	_ = writer.WriteField("query", strconv.Quote(rawQuery))
	_ = writer.WriteField("description", "profile header data")

	// Close the writer to finalize the multipart form-data
	err := writer.Close()
	if err != nil {
		slog.Error("Error closing writer", "error", err)
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", myheritageUrl, &payload)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	req.Header.Add("authorization", "Bearer "+apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response", "error", err)
		return nil, err
	}

	slog.Info("Response body", "body", string(body))

	if res.StatusCode != http.StatusOK {
		slog.Error("Non-OK HTTP status", "status", res.StatusCode, "body", string(body))
		return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	var profile ProfileHeaderResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &profile.Data.Profile, nil
}

type ProfileDetailsResponse struct {
	Data struct {
		Profile struct {
			Individual IndividualDetails `json:"individual"`
		} `json:"profile"`
	} `json:"data"`
}

type IndividualDetails struct {
	EventFacts struct {
		Data []EventFact `json:"data"`
	} `json:"event_facts"`
	FamilyGroups []FamilyGroup `json:"family_groups"`
}

const FamilyGroupTypeParent = "parent"

func (d *IndividualDetails) GetFatherId() string {
	for _, familyGroup := range d.FamilyGroups {
		if familyGroup.IsParentFamily && familyGroup.Type == FamilyGroupTypeParent &&
			familyGroup.Father != nil {
			return familyGroup.Father.Individual.Id
		}
	}

	return ""
}

func (d *IndividualDetails) GetMotherId() string {
	for _, familyGroup := range d.FamilyGroups {
		if familyGroup.IsParentFamily && familyGroup.Type == FamilyGroupTypeParent &&
			familyGroup.Mother != nil {
			return familyGroup.Mother.Individual.Id
		}
	}

	return ""
}

type EventFact struct {
	Id               string `json:"id"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	IsFamilyFact     bool   `json:"is_family_fact"`
	IsFactOfRelative bool   `json:"is_fact_of_relative"`
	Date             struct {
		Text string `json:"text"`
	} `json:"date"`
	Year              string `json:"year"`
	FormattedAge      string `json:"formatted_age"`
	FormattedPlace    string `json:"formatted_place"`
	CauseOfDeath      string `json:"cause_of_death"`
	Content           string `json:"content"`
	AdditionalContent string `json:"additional_content"`
	Individual        struct {
		Id string `json:"id"`
	} `json:"individual"`
	Relative  interface{} `json:"relative"`
	Spouse    interface{} `json:"spouse"`
	Hint      interface{} `json:"hint"`
	Citations struct {
		Data interface{} `json:"data"`
	} `json:"citations"`
	Notes struct {
		Data interface{} `json:"data"`
	} `json:"notes"`
	Media struct {
		Data interface{} `json:"data"`
	} `json:"media"`
}

type FamilyGroup struct {
	Type           string              `json:"type"`
	IsParentFamily bool                `json:"is_parent_family"`
	Father         *FamilyGroupMember  `json:"father"`
	Mother         *FamilyGroupMember  `json:"mother"`
	Siblings       []FamilyGroupMember `json:"siblings"`
	Spouse         *FamilyGroupMember  `json:"spouse"`
	Children       []FamilyGroupMember `json:"children"`
}

type FamilyGroupMember struct {
	RelationshipDescription string `json:"relationship_description"`
	RelationshipType        string `json:"relationship_type"`
	Individual              struct {
		Id                string      `json:"id"`
		Name              string      `json:"name"`
		Gender            string      `json:"gender"`
		AgeGroup          string      `json:"age_group"`
		Lifespan          string      `json:"lifespan"`
		PersonalPhoto     interface{} `json:"personal_photo"`
		LinkInProfilePage string      `json:"link_in_profile_page"`
	} `json:"individual"`
}

func GetProfileDetails(apiKey, profileId string) (*IndividualDetails, error) {
	// Create a Graphql request
	var graphqlRequest GraphqlRequest
	graphqlRequest.Query = `{profile(id:"` + profileId + `",lang:"EN"){individual{family_groups(relationship_prefix:"auto"){type is_parent_family father{...family_member_fragment}mother{...family_member_fragment}siblings(include_half_siblings:true){...family_member_fragment}spouse{...family_member_fragment}children{...family_member_fragment}}event_facts(hints:3){data{...fact_fragment}}insights{confirmed_record_matches_summary{...insight_summary_fragment}consistency_issues_summary{...insight_summary_fragment}relative_hints{...hint_fragment}}map_pins{data{...map_pin_fragment}}}site_membership{...site_membership_fragment}user{surname_research}birthday_greeting{...greeting_fragment}anniversary_greeting{...greeting_fragment}}}fragment fact_fragment on Fact{id type title is_family_fact is_fact_of_relative date{text}year formatted_age formatted_place cause_of_death content additional_content individual{id}relative{...fact_relative_fragment}spouse{...fact_relative_fragment}hint{...hint_fragment}citations{data{...citation_fragment}}notes{data{...note_fragment}}media{data{name link thumbnails(thumbnail_size:"96x96c"){url}}}}fragment citation_fragment on Citation{id page confidence event{id title}family_event{id title}date{text}formatted_text page_link{url name image}source{name smart_matching_site{id}image link}extended_citation{reference comment reason}smart_matching_individual{id name}}fragment note_fragment on Note{id type text subject body}fragment family_member_fragment on Relationship{relationship_description relationship_type individual{id name gender age_group lifespan personal_photo{...personal_photo_fragment}link_in_profile_page}}fragment personal_photo_fragment on Photo{thumbnails(thumbnail_size:"136x136c"){url}}fragment fact_relative_fragment on Individual{id name gender age_group personal_photo{...personal_photo_fragment}link_in_profile_page}fragment insight_summary_fragment on InsightSummary{type status count link is_accessible fields{id label value}}fragment hint_fragment on InsightHint{factor key modifier count first_source_name image}fragment map_pin_fragment on FactMapPin{location{name point{lat lng}bounds{north_east{lat lng}south_west{lat lng}}}facts{data{id is_fact_of_relative is_family_fact title date{text}formatted_place individual{id}relative{id name}spouse{name}}}}fragment sentence_fragment on StorySentence{text tokens{type text value link}}fragment site_membership_fragment on ProfileSiteMembership{member_id member_gender site_id site_creator_id role_sentence{...sentence_fragment}visit_sentence{...sentence_fragment}join_sentence{...sentence_fragment}request_sentence{...sentence_fragment}is_current_user_member_in_site can_user_contact_member can_user_contact_site_manager can_user_promote_member_to_site_manager can_user_demote_member_from_site_manager can_user_remind_member_to_visit can_user_change_member_email_for_remind_to_visit can_user_review_membership_request review_membership_request_link can_user_remove_member_from_site can_user_identify_member_in_tree can_user_edit_member_profile edit_member_profile_link can_user_edit_member_site_preferences edit_member_site_preferences_link can_user_edit_member_privacy_preferences edit_member_privacy_preferences_link can_user_change_member_email_and_password change_member_email_and_password_link can_user_view_member_public_profile view_member_public_profile_link can_user_associate_member_in_tree other_site_memberships{data{site_name site_link role}}}fragment greeting_fragment on ProfileGreeting{type date title label link}`
	graphqlRequest.Description = "profile details data"

	// Convert struct to JSON
	jsonData, err := json.Marshal(graphqlRequest)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", myheritageUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	client := &http.Client{}
	req.Header.Add("authorization", "Bearer "+apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response", "error", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("Non-OK HTTP status", "status", res.StatusCode, "body", string(body))
		return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	var profile ProfileDetailsResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &profile.Data.Profile.Individual, nil
}

func UpdateProfile(apiKey, profileId, name, description string) error {
	return nil
}

func DeleteProfile(apiKey, profileId string) error {
	return nil
}
