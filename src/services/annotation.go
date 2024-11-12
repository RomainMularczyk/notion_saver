package services

import (
	"github.com/google/uuid"
	"notion_saver/src/models"
	"notion_saver/src/repositories"
	"notion_saver/src/utils/data"
)

// Create an annotation and add block foreign key
func CreateAnnotation(
	notionAnnotation models.NotionBlockParagraphTextAnnotation,
	blockId uuid.UUID,
) (*models.Annotation, error) {
	annotation := data.BuildAnnotationContent(notionAnnotation, blockId)
	queryAnnotation := repositories.QueryAnnotation()
	return queryAnnotation.Create(annotation)
}
