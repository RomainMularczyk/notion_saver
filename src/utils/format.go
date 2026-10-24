package utils

import (
	"encoding/base64"
	"encoding/json"
)

/*
Formats the given data as JSON.
*/
func ToJson(data any) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}

// Decodes the given base64 string.
func FromBase64(data []byte) []byte {
	encodedData := base64.StdEncoding.EncodeToString(data)
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil
	}
	return decodedData
}
