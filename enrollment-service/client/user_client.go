package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GetUserByID(userID uint) (map[string]interface{}, error) {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://user-service:8081"
	}

	var result map[string]interface{}
	resp, err := resty.New().R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/users/%d", baseURL, userID))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("user not found")
	}
	return result, nil
}
