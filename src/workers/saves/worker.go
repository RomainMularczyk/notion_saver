package savesWorker

import (
	"log/slog"
	"notion_saver/src/utils"
)

// The PageWorker consumes messages from the 'notion.pages' queue
// and retrieves following pages from the Notion API
type SaveWorker struct {
	server *utils.Server
	logger *slog.Logger
}
