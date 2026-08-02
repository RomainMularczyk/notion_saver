package notionModels

type NotionTable struct {
	TableWidth      int  `json:"table_width"`
	HasColumnHeader bool `json:"has_column_header"`
	HawRowHeader    bool `json:"has_row_header"`
}

type NotionTableRow struct {
	Cells []NotionTableCell `json:"cells"`
}

type NotionTableCell struct {
	Type        string             `json:"type"`
	Text        NotionText         `json:"text"`
	Annotations []NotionAnnotation `json:"annotations"`
	PlainText   string             `json:"plain_text"`
	Href        string             `json:"href"`
}
