package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"notion_saver/src/database"
	"notion_saver/src/models"
)

type SaveRepository struct {
	DB *gorm.DB
}

func QuerySave() *SaveRepository {
	db := database.DatabaseConfig()
	return &SaveRepository{DB: db}
}

// Retrieve all existing saves in the database
func (sr *SaveRepository) RetrieveAll() []models.Save {
	var saves []models.Save

	if err := sr.DB.Find(&saves).Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to retrieve all saves: %v", err),
		)
	}
	return saves
}

// Insert a new save in the database
func (sr *SaveRepository) Create(save models.Save) (*models.Save, error) {
	result := sr.DB.Create(&save)
	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to create a new save: %s", err),
		)
		return nil, err
	} else {
		slog.Info(
			fmt.Sprintf("New save was inserted successfully: %s", &save.Id),
		)
		return &save, nil
	}
}

// Delete a save using its ID
func (sr *SaveRepository) DeleteOne(id string) error {
	result := sr.DB.Delete(&models.Save{}, id)

	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to delete save with Id: %s", id),
		)
		return err
	}
	return nil
}

// Retrieve the newest save available in the database
func (sr *SaveRepository) RetrieveLatest() (*models.Save, error) {
	var save models.Save
	result := sr.DB.Raw("SELECT id, MAX(last_save) FROM saves")
	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to retrieve the latest save: %s", err),
		)
		return nil, err
	} else {
		result.Scan(&save)
		fmt.Println(&save)
		slog.Info(
			fmt.Sprintf("Latest save was retrieved successfully: %s", &save.Id),
		)
		return &save, nil
	}
}
