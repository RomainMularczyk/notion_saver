package main

import (
	"notion_saver/src/database"
	"notion_saver/src/routes"
)

func main() {
	database.MigrateSchemas()
	routes.Router()
}
