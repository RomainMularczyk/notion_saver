package data

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"notion_saver/src/models"
	"strings"
)

func BuildBlocksData(notionBlocks models.NotionBlocks, blocksId uuid.UUID) (*[]models.BlockData, error) {
	var blockInBlocks []models.BlockData
	var block models.BlockData
	for _, result := range notionBlocks.Results {
		blockUuid, err := uuid.Parse(result.Id)
		if err != nil {
			slog.Error(
				fmt.Sprintf("Error occurred when parsing block UUID: %v", err),
			)
			return nil, err
		}
		if result.Type == "paragraph" {
			block = models.BlockData{
				Id:             blockUuid,
				Object:         block.Object,
				LastEditedTime: block.LastEditedTime,
				Type:           block.Type,
				PlainText:      BuildPageContentFromBlock(result),
				BlocksId:       blocksId,
				Annotation: BuildAnnotationContentFromBlockParagraph(
					result.Paragraph,
					blockUuid,
				),
			}
			blockInBlocks = append(blockInBlocks, block)
		}
	}
	return &blockInBlocks, nil
}

// Retrieve blocks content from Notion blocks
func BuildBlocksContent(
	notionBlocks models.NotionBlocks,
	pageId uuid.UUID,
) models.Blocks {
	var blocks models.Blocks
	blocksUuid := uuid.New()
	blocks = models.Blocks{
		Id:       blocksUuid,
		FullText: BuildPageContent(notionBlocks),
		PageId:   pageId,
	}
	return blocks
}

// Retrive text content from Notion block
func BuildPageContentFromBlock(block models.NotionBlock) string {
	var textContent strings.Builder
	for _, content := range block.Paragraph.RichText {
		textContent.WriteString(content.PlainText + " ")
	}
	if len(textContent.String()) == 0 {
		return textContent.String()
	}
	return textContent.String()[:len(textContent.String())-1]
}
