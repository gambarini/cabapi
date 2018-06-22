package model

type (
	Trips struct {
		Medallion string `json:"medallion"`
		Date      string `json:"date"`
		Total     int    `json:"total"`
	}
)
