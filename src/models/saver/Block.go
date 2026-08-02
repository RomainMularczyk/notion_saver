package saverModels

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// NotionSaver models
type Block struct {
	Id       uuid.UUID      `gorm:"primaryKey" json:"id"`
	FullText string         `json:"full_text"`
	Type     string         `json:"type"`
	Payload  datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	Hash     *string        `json:"hash"`
	PageId   uuid.UUID      `json:"page_id"`
	Page     Page           `gorm:"foreignKey:PageId;references:Id"`
}

// Check if two blocks have the same content
// by comparing their hashes
func (b *Block) HasSameContent(block Block) bool {
	return b.Hash == block.Hash
}

func (b *Block) ComputeHash() *string {
	hashSum := sha256.Sum256([]byte(b.Payload))
	hash := hex.EncodeToString(hashSum[:])
	return &hash
}
