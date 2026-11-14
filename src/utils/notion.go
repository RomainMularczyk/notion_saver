package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	notionModels "notion_saver/src/models/notion"

	"github.com/google/uuid"
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
	Code int
	Kind NotionErrorKind
	Err  error
}

type NotionErrorKind string

const (
	NotionRateLimitError NotionErrorKind = "rate_limit_error"
	NotionClientError    NotionErrorKind = "client_error"
	NotionServerError    NotionErrorKind = "server_error"
	NotionTransportError NotionErrorKind = "transport_error"
	NotionParsingError   NotionErrorKind = "parsing_error"
	NotionRequestError   NotionErrorKind = "request_error"
)

/*
	Check whether the NotionError is retryable or not.

For example, in cases of a client error (400), except rate limiting (429),
the request is considered as not retryable and the message should be discarded
right away.
*/
func (e *NotionError) IsRetryable() bool {
	if e == nil {
		return false
	}
	switch e.Kind {
	case NotionRateLimitError, NotionServerError, NotionTransportError:
		return true
	default:
		return false
	}
}

// Return the original error
func (e *NotionError) Unwrap() error {
	return e.Err
}

// Print the error message
func (e *NotionError) Error() string {
	if e == nil {
		return "Notion API error: <nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("Notion API error (code=%d): %v", e.Code, e.Err)
	}
	return fmt.Sprintf("Notion API error (code=%d)", e.Code)
}

type Notion struct {
	server        *Server
	logger        *slog.Logger
	endpoint      string
	headerVersion string
}

func NewNotion(server *Server) *Notion {
	return &Notion{
		server:        server,
		logger:        server.Logger.With("module", "utils/notion"),
		endpoint:      server.Env["NOTION_API_ENDPOINT"],
		headerVersion: server.Env["NOTION_API_HEADER_VERSION"],
	}
}

// Set the headers conforming to the Notion API requirements
// for more information, refer to https://developers.notion.com/reference/intro
func (n *Notion) setHeaders(request *http.Request) {
	if n.server.Env["NOTION_API_KEY"] == "" {
		panic("Could not find NOTION_API_KEY environment variable.")
	}
	if n.server.Env["NOTION_API_ENDPOINT"] == "" {
		panic("Could not find NOTION_API_ENDPOINT environment variable.")
	}
	if n.server.Env["NOTION_API_HEADER_VERSION"] == "" {
		panic("Could not find NOTION_API_HEADER_VERSION environment variable.")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Notion-Version", n.headerVersion)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", n.server.Env["NOTION_API_KEY"]))
}

// Filter Notion pages by type (either 'page' or 'database')
func (n *Notion) filterPages(
	pageType NotionFilterType,
	nextCursor *string,
) []byte {
	if !pageType.Valid() {
		errMsg := "Notion API: Invalid filter type. Should be either 'page' or 'database'."
		n.logger.Error(errMsg, slog.String("page_type", string(pageType)))
		panic(errMsg)
	}

	var filter map[string]any

	// NOTE: This is the way we should build Notion API body filters
	// for more information, refer to https://developers.notion.com/reference/post-search
	if nextCursor != nil {
		n.logger.Debug("NextCursor is defined.", slog.String("next_cursor", *nextCursor))
		filter = map[string]any{
			"filter": map[string]any{
				"value":    pageType,
				"property": "object",
			},
			"start_cursor": *nextCursor,
		}
	} else {
		n.logger.Debug("NextCursor is not defined.")
		filter = map[string]any{
			"filter": map[string]any{
				"value":    pageType,
				"property": "object",
			},
		}
	}

	body, err := json.Marshal(filter)
	if err != nil {
		errMsg := "Error when marshalling JSON."
		n.logger.Error(errMsg, slog.String("error", err.Error()))
		panic(errMsg)
	}

	return body
}

// Create the HTTP request
func (n *Notion) setRequest(
	method string,
	url string,
	body []byte,
) (*http.Request, *NotionError) {
	request, err := http.NewRequest(
		method,
		n.endpoint+url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		n.logger.Error(
			"Error when creating the request.",
			slog.Int("status_code", 0),
			slog.String("error", err.Error()),
		)
		return nil, &NotionError{
			Err:  err,
			Kind: NotionRequestError,
			Code: 0,
		}
	}

	n.setHeaders(request)

	return request, nil
}

// Send the HTTP request to Notion API
func (n *Notion) sendRequest(request *http.Request) (*http.Response, *NotionError) {
	n.logger.Debug(
		"Sending request to Notion API.",
		slog.String("url", request.URL.String()),
	)
	response, err := n.server.HTTPClient.Do(request)
	if err != nil {
		n.logger.Error(
			"Error when sending the POST request.",
			slog.Int("status_code", 0),
			slog.String("kind", string(NotionRequestError)),
			slog.String("url", request.URL.String()),
			slog.String("error", err.Error()),
		)
		return nil, &NotionError{
			Err:  err,
			Kind: NotionRequestError,
			Code: 0,
		}
	}
	n.logger.Debug(
		"Notion API response.",
		slog.String("url", request.URL.String()),
		slog.Int("status_code", response.StatusCode),
	)

	notionErr := n.filterAPIError(response)
	if notionErr != nil {
		n.logger.Error(
			"Error response returned by Notion API.",
			slog.Int("status_code", notionErr.Code),
			slog.String("kind", string(notionErr.Kind)),
		)
		response.Body.Close()
		return nil, notionErr
	}

	return response, nil
}

// Search Notion pages by type (either 'page' or 'database')
func (n *Notion) SearchPages(
	pageType NotionFilterType,
	nextCursor *string,
) (*notionModels.NotionPages, *NotionError) {
	pageFilter := n.filterPages(pageType, nextCursor)
	request, notionErr := n.setRequest("POST", "/search", pageFilter)
	if notionErr != nil {
		return nil, notionErr
	}

	response, notionErr := n.sendRequest(request)
	if notionErr != nil {
		return nil, notionErr
	}
	defer response.Body.Close()

	var notionPages notionModels.NotionPages
	err := json.NewDecoder(response.Body).Decode(&notionPages)
	if err != nil {
		n.logger.Error(
			"Error when unmarshalling Notion Pages.",
			slog.String("error", err.Error()),
			slog.String("url", request.URL.String()),
		)
		return nil, &NotionError{
			Err:  err,
			Kind: NotionParsingError,
			Code: response.StatusCode,
		}
	}
	n.logger.Info("Notion API responded successfully.")

	return &notionPages, nil
}

// Get the blocks of a Notion page
func (n *Notion) GetPageBlocks(
	pageId uuid.UUID,
	nextCursor *string,
) (*notionModels.NotionBlocks, *NotionError) {
	path := "/blocks/" + pageId.String() + "/children"

	if nextCursor != nil {
		query := url.Values{}
		query.Set("start_cursor", *nextCursor)
		path += "?" + query.Encode()
	}
	request, notionErr := n.setRequest("GET", path, nil)
	if notionErr != nil {
		return nil, notionErr
	}

	response, notionErr := n.sendRequest(request)
	if notionErr != nil {
		return nil, notionErr
	}
	defer response.Body.Close()

	var blocks notionModels.NotionBlocks
	err := json.NewDecoder(response.Body).Decode(&blocks)
	if err != nil {
		n.logger.Error(
			"Error when unmarshalling Notion Blocks.",
			slog.String("error", err.Error()),
			slog.String("url", request.URL.String()),
		)
		return nil, &NotionError{
			Err:  err,
			Kind: NotionParsingError,
			Code: response.StatusCode,
		}
	}

	return &blocks, nil
}

// Filter Notion API errors to enable retry behaviors
func (n *Notion) filterAPIError(response *http.Response) *NotionError {
	if response.StatusCode == http.StatusTooManyRequests {
		return &NotionError{
			Kind: NotionRateLimitError,
			Code: response.StatusCode,
		}
	}
	if response.StatusCode >= 500 {
		return &NotionError{
			Kind: NotionServerError,
			Code: response.StatusCode,
		}
	}
	if response.StatusCode >= 400 {
		return &NotionError{
			Kind: NotionClientError,
			Code: response.StatusCode,
		}
	}
	return nil
}
