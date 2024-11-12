package data

import (
	"github.com/google/uuid"
	"notion_saver/src/models"
)

// Retrieve annotation content from Notion Annotation
func BuildAnnotationContent(
	notionAnnotation models.NotionBlockParagraphTextAnnotation,
	blockId uuid.UUID,
) models.Annotation {
	var annotation models.Annotation
	annotationUuid := uuid.New()
	annotation = models.Annotation{
		Id:            annotationUuid,
		Bold:          notionAnnotation.Bold,
		Italic:        notionAnnotation.Italic,
		Strikethrough: notionAnnotation.Strikethrough,
		Underline:     notionAnnotation.Underline,
		Code:          notionAnnotation.Code,
		Color:         notionAnnotation.Color,
		BlockId:       blockId,
	}
	return annotation
}

// Retrieve annotation content from Notion block paragraph
func BuildAnnotationContentFromBlockParagraph(
	paragraph models.NotionBlockParagraph,
	blockId uuid.UUID,
) []models.Annotation {
	var annotations []models.Annotation
	for _, section := range paragraph.RichText {
		annotation := BuildAnnotationContent(section.Annotations, blockId)
		annotations = append(annotations, annotation)
	}
	return annotations
}
