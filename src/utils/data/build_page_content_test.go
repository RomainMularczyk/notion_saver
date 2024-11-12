package data

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"notion_saver/src/models"
	"testing"
)

// Verify that the content of a Notion page is properly gathered into its entire
// plain text content
func TestBuildPageContent(t *testing.T) {
	notionBlocksAsJson := `{
		"results": [
			{
				"object": "block",
				"id": "block-1",
				"last_edited_time": "2022-11-06T17:02:00.000Z",
				"has_children": false,
				"type": "paragraph",
				"paragraph": {
					"rich_text": [
						{
							"type": "text",
							"annotations": {
								"bold": true,
								"italic": false,
								"strikethrough": false,
								"underline": false,
								"code": false,
								"color": "default"
							},
							"plain_text": "This is the first paragraph."
						}
					]
				}
			},
			{
				"object": "block",
				"id": "block-2",
				"last_edited_time": "2022-11-06T18:02:00.000Z",
				"has_children": false,
				"type": "paragraph",
				"paragraph": {
					"rich_text": [
						{
							"type": "text",
							"annotations": {
								"bold": false,
								"italic": true,
								"strikethrough": false,
								"underline": true,
								"code": false,
								"color": "default"
							},
							"plain_text": "This is the second paragraph with italic and underline."
						}
					]
				}
			}
		]
	}`
	expected := "This is the first paragraph. This is the second paragraph with italic and underline."
	var blocks models.NotionBlocks
	json.Unmarshal([]byte(notionBlocksAsJson), &blocks)
	result := BuildPageContent(blocks)
	assert.Equal(t, result, expected)
}
