package cmd

import (
	"fmt"
	"log"
	"errors"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"io/ioutil"

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

func getCategory(impl string) ([]Category, error) {
	var categories []Category
	switch impl {
	case "local":

		data, err := ioutil.ReadFile("./config/default.yml")
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &categories)
		if err != nil {
			return nil, err
		}

	case "chatgpt":
		
		apiKey := "your-api-key"
		client := resty.New()
		response, err := client.R().
			SetAuthToken(apiKey).
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"model": "gpt-3.5-turbo",
				"messages": []interface{}{map[string]interface{}{
					"role":"system",
					"content": "Hi can you tell me what is the factorial of 10?"}},
				"max_tokens": 50,
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

		content := data//["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		fmt.Println(content)


	default:
		return nil, errors.New(fmt.Sprintf(" Error: Invalide Implementation Type."))

	}
	return categories, nil
}

