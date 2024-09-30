package myheritage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateProfile(apiKey, name, description string) (string, error) {
	url := "https://api.myheritage.com/family-trees" // Replace with actual MyHeritage API endpoint

	familyTree := Profile{
		Name:        name,
		Description: description,
	}

	jsonBody, err := json.Marshal(familyTree)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey) // Replace with actual API key

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var createdFamilyTree Profile
	err = json.Unmarshal(body, &createdFamilyTree)
	if err != nil {
		return "", err
	}

	return createdFamilyTree.ID, nil
}

func GetProfile(apiKey, familyTreeID string) (*Profile, error) {
	url := fmt.Sprintf("https://api.myheritage.com/family-trees/%s", familyTreeID) // Replace with the actual MyHeritage API endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey) // Replace with actual API key

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get family tree: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var familyTree Profile
	err = json.Unmarshal(body, &familyTree)
	if err != nil {
		return nil, err
	}

	return &familyTree, nil
}

func UpdateProfile(apiKey, familyTreeID, name, description string) error {
	url := fmt.Sprintf("https://api.myheritage.com/family-trees/%s", familyTreeID) // Replace with actual MyHeritage API endpoint

	familyTree := Profile{
		Name:        name,
		Description: description,
	}

	jsonBody, err := json.Marshal(familyTree)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey) // Replace with actual API key

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update family tree: %s", resp.Status)
	}

	return nil
}

func DeleteProfile(apiKey, familyTreeID string) error {
	url := fmt.Sprintf("https://api.myheritage.com/family-trees/%s", familyTreeID) // Replace with actual MyHeritage API endpoint

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey) // Replace with actual API key

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete family tree: %s", resp.Status)
	}

	return nil
}
