package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
	"notion_saver/src/database"
	"notion_saver/src/models"
	"strconv"
)

type PageRepository struct {
	DB *gorm.DB
}

func QueryPage() *PageRepository {
	db := database.DatabaseConfig()
	return &PageRepository{DB: db}
}

// Retrieve all existing pages in the database
func (pr *PageRepository) RetrieveAll() []models.Page {
	var pages []models.Page

	if err := pr.DB.Find(&pages).Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to retrieve all pages: %v", err),
		)
	}
	return pages
}

// Get a page using its ID
func (pr *PageRepository) RetrieveOne(id string) models.Page {
	var page models.Page

	uuid, err := uuid.Parse(id)
	if err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to parse the ID: %s", id),
		)
	}

	result := pr.DB.Model(models.Page{Id: uuid}).First(&page)

	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to retrieve page with Id: %s", id),
		)
	}
	return page
}

// Delete a page using its ID
func (pr *PageRepository) DeleteOne(id string) error {
	result := pr.DB.Delete(&models.Page{}, id)

	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to delete page with Id: %s", id),
		)
		return err
	}
	return nil
}

// Insert a new page in the database
func (pr *PageRepository) Create(page models.Page) (*models.Page, error) {
	result := pr.DB.Create(&page)
	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to create a new page: %s", err),
		)
		return nil, err
	} else {
		slog.Info(
			fmt.Sprintf("New page was inserted successfully: %s", page.Id),
		)
		return &page, nil
	}
}

// Insert many new pages in the database
func (pr *PageRepository) CreateMany(pages []models.Page) (*[]models.Page, error) {
	for i, page := range pages {
		_, err := pr.Create(page)
		if err != nil {
			slog.Error(
				fmt.Sprintf(
					"An error occurred at iteration %s (over %s iterations).",
					strconv.Itoa(i),
					strconv.Itoa(len(pages)),
				),
			)
			return nil, err
		}
	}

	return &pages, nil
}
