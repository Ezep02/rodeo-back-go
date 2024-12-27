package utils

import (
	"log"
	"time"
)

func TimeParser(date string) (*time.Time, error) {

	layout := "01-02-2006"
	parsedDate, err := time.Parse(layout, date)

	if err != nil {
		log.Println("Error:", err)
		return &time.Time{}, err
	}

	return &parsedDate, nil
}
