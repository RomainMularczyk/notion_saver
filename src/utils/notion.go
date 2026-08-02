package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	notionModels "notion_saver/src/models/notion"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type NotionFilterType string

const (
	FilterTypePage     NotionFilterType = "page"
	FilterTypeDatabase NotionFilterType = "database"
)

// Check if a NotionFilterType is valid
func (t NotionFilterType) Valid() bool {
	return t == FilterTypePage || t == FilterTypeDatabase
}

type NotionError struct {
	Code          int
	IsRateLimited bool
	IsClientError bool
	IsServerError bool
	IsMalformed   bool
}

func (e NotionError) Error() string {
	return fmt.Sprintf(
		"An error occurred when retrieving data from Notion API. Code: %v",
		e.Code,
	)
}

type Notion struct {
	server        *Server
	logger        *slog.Logger
	Endpoint      string
	HeaderVersion string
}

func NewNotion(server *Server) *Notion {
	return &Notion{
		server:        server,
		logger:        server.Logger.With("module", "utils/notion"),
		Endpoint:      "https://api.notion.com/v1",
		HeaderVersion: "2026-03-11",
	}
}

// Set the headers conforming to the Notion API requirements
// for more information, refer to https://developers.notion.com/reference/intro
func (n *Notion) setHeaders(request *http.Request) {
	if err := godotenv.Load(".env"); err != nil {
		errMsg := fmt.Sprintf("Could not find NOTION_API_BEARER_TOKEN environment variable.")
		n.logger.Error(errMsg)
		panic(errMsg)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Notion-Version", n.HeaderVersion)
	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", n.server.Env["NOTION_API_KEY"]),
	)
}

// Filter Notion pages by type (either 'page' or 'database')
func (n *Notion) filterPages(
	pageType NotionFilterType,
	nextCursor *string,
) []byte {
	if !pageType.Valid() {
		errMsg := fmt.Sprintf(
			"Notion API: Invalid filter type. Should be either 'page' or 'database'. Received: %v", pageType,
		)
		n.logger.Error(errMsg)
		panic(errMsg)
	}

	var filter map[string]any

	// This is the way we should build Notion API body filters
	// for more information, refer to https://developers.notion.com/reference/post-search
	if nextCursor != nil {
		n.logger.Debug("NextCursor is defined", slog.String("nextCursor", *nextCursor))
		filter = map[string]any{
			"filter": map[string]any{
				"value":    pageType,
				"property": "object",
			},
			"start_cursor": *nextCursor,
		}
	} else {
		n.logger.Debug("NextCursor is not defined")
		filter = map[string]any{
			"filter": map[string]any{
				"value":    pageType,
				"property": "object",
			},
		}
	}

	body, err := json.Marshal(filter)
	if err != nil {
		errMsg := fmt.Sprintf("Error when marshalling JSON: %v", err)
		n.logger.Error(errMsg)
		panic(errMsg)
	}

	return body
}

// Create the HTTP request
func (n *Notion) setRequest(
	method string,
	url string,
	body []byte,
) (*http.Request, error) {
	request, err := http.NewRequest(
		method,
		n.Endpoint+url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		n.logger.Error(fmt.Sprintf("Error when creating a request: %v", err))
		return nil, err
	}

	n.setHeaders(request)

	return request, nil
}

// Search Notion pages by type (either 'page' or 'database')
func (n *Notion) SearchPages(
	pageType NotionFilterType,
	nextCursor *string,
) (*notionModels.NotionPages, error) {
	request, err := n.setRequest(
		"POST",
		"/search",
		n.filterPages(pageType, nextCursor),
	)
	if err != nil {
		n.logger.Error(fmt.Sprintf("Error when creating a request: %v", err))
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		n.logger.Error(
			fmt.Sprintf("Error when sending the POST request: %v", err),
			slog.String("URL", request.URL.String()),
		)
		return nil, err
	}

	defer response.Body.Close()
	var notionPages notionModels.NotionPages

	err = json.NewDecoder(response.Body).Decode(&notionPages)

	if err != nil {
		n.logger.Error(
			fmt.Sprintf("Error when unmarshalling Notion Pages: %v", err),
			slog.String("URL", request.URL.String()),
		)
		return nil, err
	}
	n.logger.Info("Notion API responsed successfully.")

	return &notionPages, nil
}

// Get the blocks of a Notion page
func (n *Notion) GetPageBlocks(
	pageId uuid.UUID,
	nextCursor *string,
) (*notionModels.NotionBlocks, error) {
	var request *http.Request
	var err error

	if nextCursor != nil {
		request, err = n.setRequest(
			"GET",
			"/blocks/"+pageId.String()+"/children?start_cursor="+*nextCursor,
			nil,
		)
	} else {
		request, err = n.setRequest(
			"GET",
			"/blocks/"+pageId.String()+"/children",
			nil,
		)
	}
	if err != nil {
		n.logger.Error(fmt.Sprintf("Error when creating a request: %v", err))
		return nil, err
	}

	client := &http.Client{}
	n.logger.Debug(
		"Sending request to Notion API.",
		slog.String("URL", request.URL.String()),
	)
	response, err := client.Do(request)
	n.logger.Debug(
		"Notion API response.",
		slog.Int("Status Code", response.StatusCode),
	)
	if err != nil {
		n.logger.Error(
			fmt.Sprintf("Error when sending the POST request: %v", err),
			slog.String("URL", request.URL.String()),
		)
		return nil, err
	}
	notionErr := n.filterAPIError(response)
	if notionErr != nil {
		return nil, notionErr
	}

	defer response.Body.Close()
	var blocks notionModels.NotionBlocks

	err = json.NewDecoder(response.Body).Decode(&blocks)
	if err != nil {
		n.logger.Error(
			fmt.Sprintf("Error when unmarshalling Notion Blocks: %v", err),
			slog.String("URL", request.URL.String()),
		)
		return nil, &NotionError{
			IsRateLimited: false,
			IsClientError: false,
			IsServerError: false,
			IsMalformed:   true,
			Code:          response.StatusCode,
		}
	}

	return &blocks, nil
}

// Filter Notion API errors to enable retry behaviors
func (n *Notion) filterAPIError(response *http.Response) *NotionError {
	if response.StatusCode == 429 {
		return &NotionError{
			IsRateLimited: true,
			IsClientError: false,
			IsServerError: false,
			Code:          response.StatusCode,
		}
	}
	if response.StatusCode >= 400 {
		return &NotionError{
			IsRateLimited: false,
			IsClientError: true,
			IsServerError: false,
			Code:          response.StatusCode,
		}
	}
	if response.StatusCode >= 500 {
		return &NotionError{
			IsRateLimited: false,
			IsClientError: false,
			IsServerError: true,
			Code:          response.StatusCode,
		}
	}
	return nil
}
