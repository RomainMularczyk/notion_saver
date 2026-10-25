package blocksWorkerModels

import "github.com/google/uuid"

type BlockMessage struct {
	NextCursor      *string   `json:"next_cursor"`
	PageId          uuid.UUID `json:"page_id"`
	PaginationIndex int       `json:"pagination_index"`
	RetryCount      int       `json:"retry_count"`
}
