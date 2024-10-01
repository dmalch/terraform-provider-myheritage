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
	Data ProfileData `json:"data"`
}

type ProfileData struct {
	Profile Profile `json:"profile"`
}

type Profile struct {
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Individual Individual `json:"individual"`
}

type Individual struct {
	ID        string `json:"id"`
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

const profileHeaderUrl = "https://familygraphql.myheritage.com/profile_header_data/"

func GetProfile(apiKey, profileId string) (*Profile, error) {
	// Create a buffer to hold the multipart form-data
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)

	_ = writer.WriteField("bearer_token", apiKey)
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
	req, err := http.NewRequest("POST", profileHeaderUrl, &payload)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

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

func UpdateProfile(apiKey, familyTreeID, name, description string) error {
	return nil
}

func DeleteProfile(apiKey, familyTreeID string) error {
	return nil
}
