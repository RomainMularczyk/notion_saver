package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Verify that the hashing function provides the right corresponding
// hashsum for the given text content
func TestHashPage(t *testing.T) {
	text := "New page"
	result := HashPage(text)
	expected := "40126b0e9d2e5f2f3a04b635f449a480c469b2ab413e0535b2cfce2865ac8644"
	assert.Equal(t, result, expected)
}

// Verify that the hashing function provides a different hash value
// on a slight string change
func TestHashPageNotEqual(t *testing.T) {
	text := "A new page"
	result := HashPage(text)
	expected := "40126b0e9d2e5f2f3a04b635f449a480c469b2ab413e0535b2cfce2865ac8644"
	assert.NotEqual(t, result, expected)
}

// Verify that comparing two hashes of two different strings returns
// "false" when compared with the IsSamePage method
func TestIsSamePage(t *testing.T) {
	textA := "A new page"
	textB := "An old page"
	result := IsSamePage(textA, textB)
	assert.False(t, result)
}
