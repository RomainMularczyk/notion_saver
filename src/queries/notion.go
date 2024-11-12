package queries

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log/slog"
	"net/http"
	"notion_saver/src/models"
	"os"
)

// Retrieve all pages from the Notion API.
func GetAllPages() models.NotionPages {
	err := godotenv.Load(".env")
	url := "https://api.notion.com/v1"
	request, err := http.NewRequest(
		"POST",
		url+"/search",
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when creating a request: %v", err))
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Notion-Version", "2022-06-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NOTION_API_KEY")))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when sending the POST request: %v", err))
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when reading response body: %s", err))
	}
	var pages models.NotionPages
	if err := json.Unmarshal(body, &pages); err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalling JSON: %s", err))
	}

	if response.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Notion API answered with error: %v", err))
	}

	return pages
}

// Retrieve all blocks from page ID.
func GetNotionBlocks(id uuid.UUID) models.NotionBlocks {
	err := godotenv.Load(".env")
	url := "https://api.notion.com/v1"
	request, err := http.NewRequest(
		"GET",
		url+"/blocks/"+id.String()+"/children",
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when creating a request: %v", err))
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Notion-Version", "2022-06-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NOTION_API_KEY")))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when sending the POST request: %v", err))
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when reading response body: %s", err))
	}
	var pageContent models.NotionBlocks
	if err := json.Unmarshal(body, &pageContent); err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalliong JSON: %s", err))
	}

	if response.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Notion API answered with error: %v", err))
	}

	return pageContent
}
