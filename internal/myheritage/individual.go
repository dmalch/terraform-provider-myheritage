package myheritage

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
)

type IndividualResponse struct {
	Data struct {
		Individual IndividualDetails `json:"individual"`
	} `json:"data"`
}

func GetIndividual(apiKey, individualId string) (*IndividualDetails, error) {
	// Create a buffer to hold the multipart form-data
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)

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

	req, err := http.NewRequest("POST", myheritageUrl, &payload)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	req.Header.Add("authorization", "Bearer "+apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", writer.FormDataContentType())

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var individual IndividualResponse
	err = json.Unmarshal(body, &individual)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &individual.Data.Individual, nil
}

type IndividualBiographyResponse struct {
	Data struct {
		Individual IndividualBiography `json:"individual"`
	} `json:"data"`
}

type IndividualBiography struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Notes struct {
		Data []Note `json:"data"`
	} `json:"notes"`
}

type Note struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	Body    string `json:"body"`
}

func GetIndividualBiography(apiKey, individualId string) (*IndividualBiography, error) {
	// Create a Graphql request
	var graphqlRequest GraphqlRequest
	graphqlRequest.Query = `{individual(id:"` + individualId +
		`",lang:"EN"){id name is_applicable_for_biography ai_biography{data{...ai_biography_fragment}}life_story{...life_story_fragment}can_generate_live_story,notes{data{...note_fragment}}comments{data{...comment_fragment}}}}fragment ai_biography_fragment on AiBiography{id status biography_item{id url}}fragment life_story_fragment on StoryChapter{text paragraphs{sentences{...sentence_fragment}}}fragment note_fragment on Note{id type text subject body}fragment comment_fragment on Comment{id text body submitter{id name}}fragment sentence_fragment on StorySentence{text tokens{type text value link}}`
	graphqlRequest.Description = "profile biography data"

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

	req.Header.Add("authorization", "Bearer "+apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var individual IndividualBiographyResponse
	err = json.Unmarshal(body, &individual)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &individual.Data.Individual, nil
}
