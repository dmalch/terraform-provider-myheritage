package myheritage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FamilyTree struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateFamilyTree(apiKey, name, description string) (string, error) {
	url := "https://api.myheritage.com/family-trees" // Replace with actual MyHeritage API endpoint

	familyTree := FamilyTree{
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

	var createdFamilyTree FamilyTree
	err = json.Unmarshal(body, &createdFamilyTree)
	if err != nil {
		return "", err
	}

	return createdFamilyTree.ID, nil
}

// GetFamilyTree fetches the details of a family tree from MyHeritage API
func GetFamilyTree(apiKey, familyTreeID string) (*FamilyTree, error) {
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

	var familyTree FamilyTree
	err = json.Unmarshal(body, &familyTree)
	if err != nil {
		return nil, err
	}

	return &familyTree, nil
}

func UpdateFamilyTree(apiKey, familyTreeID, name, description string) error {
	url := fmt.Sprintf("https://api.myheritage.com/family-trees/%s", familyTreeID) // Replace with actual MyHeritage API endpoint

	familyTree := FamilyTree{
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

func DeleteFamilyTree(apiKey, familyTreeID string) error {
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
