package notionModels

import (
	"encoding/json"
	"fmt"
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Notion API models
type NotionBlocks struct {
	Object     string        `json:"object"`
	Results    []NotionBlock `json:"results"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor"`
	Type       string        `json:"type"`
}

type NotionBlock struct {
	Object           string                  `json:"object"`
	Id               uuid.UUID               `json:"id"`
	CreatedTime      time.Time               `json:"created_time"`
	LastEditedTime   time.Time               `json:"last_edited_time"`
	HasChildren      bool                    `json:"has_children"`
	InTrash          bool                    `json:"in_trash"`
	Type             string                  `json:"type"`
	Paragraph        *NotionBlockParagraph   `json:"paragraph,omitempty"`
	BulletedListItem *NotionBulletedListItem `json:"bulleted_list_item,omitempty"`
	NumberedListItem *NotionNumberedListItem `json:"numbered_list_item,omitempty"`
	Heading1         *NotionHeading          `json:"heading_1,omitempty"`
	Heading2         *NotionHeading          `json:"heading_2,omitempty"`
	Heading3         *NotionHeading          `json:"heading_3,omitempty"`
	Heading4         *NotionHeading          `json:"heading_4,omitempty"`
	Code             *NotionCode             `json:"code,omitempty"`
	ChildPage        *NotionChildPage        `json:"child_page,omitempty"`
	ChildDatabase    *NotionChildDatabase    `json:"child_database,omitempty"`
	Image            *NotionImageCaption     `json:"image,omitempty"`
	LinkPreview      *NotionLinkPreview      `json:"link_preview,omitempty"`
	Equation         *NotionEquation         `json:"equation,omitempty"`
	Embed            *NotionEmbed            `json:"embed,omitempty"`
	Callout          *NotionCallout          `json:"callout,omitempty"`
	Toggle           *NotionToggle           `json:"toggle,omitempty"`
	Table            *NotionTable            `json:"table,omitempty"`
	TableRow         *NotionTableRow         `json:"table_row,omitempty"`
	Quote            *NotionQuote            `json:"quote,omitempty"`
}

// Convert a Notion API Blocks into a NotionSaver Blocks
func (b *NotionBlocks) ToSaverFormat(pageId uuid.UUID) []saverModels.Block {
	newBlocks := make([]saverModels.Block, 0)

	for _, block := range b.Results {
		payloadBytes, err := json.Marshal(block)
		if err != nil {
			slog.Error(fmt.Sprintf("Error marshalling block: %s", err))
			continue
		}
		payload := datatypes.JSON(payloadBytes)
		fulltext := block.GetFullText()

		newBlock := saverModels.Block{
			Id:       block.Id,
			FullText: fulltext,
			Type:     block.Type,
			Payload:  payload,
			PageId:   pageId,
		}
		newBlock.Hash = newBlock.ComputeHash()

		newBlocks = append(newBlocks, newBlock)
	}

	return newBlocks
}

// Retrieve the full text content of a Notion block
func (b *NotionBlock) GetFullText() string {
	if b.Type == "paragraph" {
		return b.paragraphFullText()
	}
	if b.Type == "bulleted_list_item" {
		return b.bulletedListItemFullText()
	}
	if b.Type == "numbered_list_item" {
		return b.numberedListItemFullText()
	}
	if b.Type == "heading_1" || b.Type == "heading_2" || b.Type == "heading_3" || b.Type == "heading_4" {
		return b.headingFullText()
	}
	if b.Type == "code" {
		return b.codeFullText()
	}
	if b.Type == "equation" {
		return b.equationFullText()
	}
	if b.Type == "child_page" {
		return ""
	}
	if b.Type == "child_database" {
		return ""
	}
	if b.Type == "image" {
		return ""
	}
	if b.Type == "link_preview" {
		return b.linkPreviewFullText()
	}
	if b.Type == "embed" {
		return b.embedFullText()
	}
	if b.Type == "callout" {
		return b.calloutFullText()
	}
	if b.Type == "toggle" {
		return b.toggleFullText()
	}
	if b.Type == "table" {
		return ""
	}
	if b.Type == "table_row" {
		return b.tableRowFullText()
	}
	if b.Type == "quote" {
		return b.quoteFullText()
	}

	panic(fmt.Sprintf("Type %s not supported", b.Type))
}

// Retrieve the full text content of a Notion quote
func (b *NotionBlock) quoteFullText() string {
	return b.concatenateRichText(b.Quote.RichText)
}

// Retrieve the full text content of a Notion table row
func (b *NotionBlock) tableRowFullText() string {
	var textContent strings.Builder
	for _, cell := range b.TableRow.Cells {
		textContent.WriteString(cell.PlainText + " | ")
	}
	if len(textContent.String()) == 0 {
		return textContent.String()
	}
	return textContent.String()[:len(textContent.String())-1]
}

// Retrieve the full text content of a Notion toggle
func (b *NotionBlock) toggleFullText() string {
	return b.concatenateRichText(b.Toggle.RichText)
}

// Retrieve the full text content of a Notion callout
func (b *NotionBlock) calloutFullText() string {
	return b.concatenateRichText(b.Callout.RichText)
}

// Retrieve the full text content of a Notion embed
func (b *NotionBlock) embedFullText() string {
	return b.Embed.Url
}

// Retrieve the full text content of a Notion link preview
func (b *NotionBlock) linkPreviewFullText() string {
	return b.LinkPreview.Url
}

// Retrieve the full text content of a Notion equation
func (b *NotionBlock) equationFullText() string {
	return b.Equation.Expression
}

// Retrieve the full text content of a Notion paragraph
func (b *NotionBlock) paragraphFullText() string {
	return b.concatenateRichText(b.Paragraph.RichText)
}

// Retrieve the full text content of a Notion bulleted list item
func (b *NotionBlock) bulletedListItemFullText() string {
	return b.concatenateRichText(b.BulletedListItem.RichText)
}

// Retrieve the full text content of a Notion numbered list item
func (b *NotionBlock) numberedListItemFullText() string {
	return b.concatenateRichText(b.NumberedListItem.RichText)
}

// Retrieve the full text content of a Notion heading
func (b *NotionBlock) headingFullText() string {
	if b.Type == "heading_1" {
		return b.concatenateRichText(b.Heading1.RichText)
	}
	if b.Type == "heading_2" {
		return b.concatenateRichText(b.Heading2.RichText)
	}
	if b.Type == "heading_3" {
		return b.concatenateRichText(b.Heading3.RichText)
	}
	if b.Type == "heading_4" {
		return b.concatenateRichText(b.Heading4.RichText)
	}
	panic("Heading type not supported")
}

// Retrieve the full text content of a Notion code block
func (b *NotionBlock) codeFullText() string {
	return b.concatenateRichText(b.Code.RichText)
}

// Concatenate a list of rich text content
func (b *NotionBlock) concatenateRichText(richText []NotionRichText) string {
	var textContent strings.Builder
	for _, content := range richText {
		textContent.WriteString(content.PlainText + " ")
	}
	if len(textContent.String()) == 0 {
		return textContent.String()
	}
	return textContent.String()[:len(textContent.String())-1]
}
