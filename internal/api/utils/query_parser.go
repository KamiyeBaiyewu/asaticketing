package utils

import (
	"net/url"
	"time"
)

// TimeParam - get the time value from the request
func TimeParam(query url.Values, name string) (time.Time, error) {

	// returns now as defult time

	t:= time.Now()
	value := query.Get(name)
	if value == ""{
		return t, nil
	}

	parsed, err := time.Parse(time.RFC3339,value)
	if err != nil {
		return t, err
	}

	return parsed, nil
}
