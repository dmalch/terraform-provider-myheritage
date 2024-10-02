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

const individualUrl = "https://familygraphql.myheritage.com/individual_data_with_hints_query/"

type IndividualResponse struct {
	Data struct {
		Individual IndividualDetails `json:"individual"`
	} `json:"data"`
}

func GetIndividual(apiKey, individualId string) (*IndividualDetails, error) {
	// Create a buffer to hold the multipart form-data
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)

	_ = writer.WriteField("bearer_token", apiKey)
	// Add the query part
	rawQuery := `query($individualId:String!,$lang:String,$showHintsSetting:BigInt!,$relationshipPrefix:String!){individual(id:$individualId,lang:$lang){can_edit,link_in_profile_page,insights(hints:$showHintsSetting){first_name_hint{...hint_fragment}personal_photo_hint{...hint_fragment}relative_hints{...hint_fragment}}family_groups(relationship_prefix:$relationshipPrefix){type is_parent_family father{...family_member_fragment}mother{...family_member_fragment}siblings(include_half_siblings:true){...family_member_fragment}spouse{...family_member_fragment}children{...family_member_fragment}}event_facts(hints:$showHintsSetting){data{hint{...hint_fragment}...fact_fragment}}}}fragment fact_fragment on Fact{id type title is_family_fact is_fact_of_relative date{text}year formatted_age formatted_place cause_of_death content additional_content individual{id}relative{...fact_relative_fragment}spouse{...fact_relative_fragment}hint{...hint_fragment}citations{data{...citation_fragment}}notes{data{...note_fragment}}media{data{name link thumbnails(thumbnail_size:"96x96c"){url}}}address{country_code}}fragment citation_fragment on Citation{id page confidence event{id title}family_event{id title}date{text}formatted_text page_link{url name image}source{name smart_matching_site{id}image link}extended_citation{reference comment reason}smart_matching_individual{id name}}fragment note_fragment on Note{id type text subject body}fragment fact_relative_fragment on Individual{id name gender age_group personal_photo{...personal_photo_fragment}link_in_profile_page}fragment hint_fragment on InsightHint{factor key modifier count first_source_name image}fragment personal_photo_fragment on Photo{thumbnails(thumbnail_size:"136x136c"){url}}fragment family_member_fragment on Relationship{relationship_description relationship_type individual{id name gender age_group lifespan personal_photo{...personal_photo_fragment}link_in_profile_page}}`
	_ = writer.WriteField("query", strconv.Quote(rawQuery))
	_ = writer.WriteField("variables", `{"individualId":"`+individualId+`","lang":"EN","showHintsSetting":3,"relationshipPrefix":"auto"}`)
	_ = writer.WriteField("description", "individual data with hints query")

	// Close the writer to finalize the multipart form-data
	err := writer.Close()
	if err != nil {
		slog.Error("Error closing writer", "error", err)
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", individualUrl, &payload)
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

	var individual IndividualResponse
	err = json.Unmarshal(body, &individual)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &individual.Data.Individual, nil
}
