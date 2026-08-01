package pagesWorker

import (
	"log/slog"
	"notion_saver/src/utils"
)

// The PageWorker consumes messages from the 'notion.saves' queue
// and retrieves following saves from the Notion API
type PageWorker struct {
	server *utils.Server
	logger *slog.Logger
}
