package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"notion_saver/src/database"
	"notion_saver/src/models"
)

// Insert a new annotation in the database
type AnnotationRepository struct {
	DB *gorm.DB
}

func QueryAnnotation() *AnnotationRepository {
	db := database.DatabaseConfig()
	return &AnnotationRepository{DB: db}
}

// Insert a new annotation in the database
func (ar *AnnotationRepository) Create(
	annotation models.Annotation,
) (*models.Annotation, error) {
	result := ar.DB.Create(&annotation)
	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf(
				"An error occurred when trying to create a new annotation: %v",
				err,
			),
		)
		return nil, err
	} else {
		slog.Info(
			fmt.Sprintf(
				"New annotation was inserted successfully: %s",
				annotation.Id,
			),
		)
		return &annotation, nil
	}
}
