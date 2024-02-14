package config

import (
	"os"
	"fmt"
	"log"
	"errors"
	"gopkg.in/yaml.v2"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

const (
    apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

type Question struct {
	Text string `yaml:"question"`
}

type Category struct {
	Name string `yaml:"category"`
	Questions []Question `yaml:"questions"`
}

func GetCategory(impl string) ([]Category, error) {
	var categories []Category
	switch impl {
	case "local":

		data, err := os.ReadFile("./config/default.yml")
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &categories)
		if err != nil {
			return nil, err
		}

	case "chatgpt":
		
		apiKey := "replace-me"
		client := resty.New()
		response, err := client.R().
			SetAuthToken(apiKey).
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"model": "gpt-3.5-turbo",
				"messages": []interface{}{map[string]interface{}{
					"role":"system",
					"content": "You are an interviewer for an AWS Cloud DevOps engineer. You need to create 3 categories of questions with 3 questions each that can be ask to an candidate. The return format of the category and related questions needs to be in a yaml format. Do not include the answer, do not add any number after the field name 'question' and each questions are part of a list, and each category are part of a list."}},
				"max_tokens": 500,
			}).
			Post(apiEndpoint)

		if err != nil {
			log.Fatalf("Error while sending send the request: %v", err)
		}

		body := response.Body()

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error while decoding JSON response:", err)
		}

		content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		//fmt.Println(content)

		err = yaml.Unmarshal([]byte(content), &categories)
		if err != nil {
			return nil, err
		}


	default:
		return nil, errors.New("error invalide implementation type")

	}
	return categories, nil
}

