package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GetCourseByID(courseID uint) (map[string]interface{}, error) {
	baseURL := os.Getenv("COURSE_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://course-service:8082"
	}

	var result map[string]interface{}
	resp, err := resty.New().R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/courses/%d", baseURL, courseID))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("course not found")
	}
	return result, nil
}
