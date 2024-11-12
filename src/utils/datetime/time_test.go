package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Verify that expected UTC format datetime is correctly parsed
func TestParseTime(t *testing.T) {
	datetime := "2020-11-20T20:30:56.000Z"
	result := Parse(datetime)
	expected := time.Date(
		2020,
		11,
		20,
		20,
		30,
		56,
		0,
		time.UTC,
	)
	assert.Equal(t, result, expected)
}
