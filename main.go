package main

import (
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "encoding/json"
		"bytes"
		"io/ioutil"
		"net/http"
		"os"
		"fmt"
)

func HandleRequest(ctx context.Context, event json.RawMessage) (json.RawMessage, error) {
	url := os.Getenv("BASE_URL")

	var f interface{}
	_  = json.Unmarshal(event, &f)
	m := f.(map[string]interface{})["directive"].(map[string]interface{})["endpoint"].(map[string]interface{})["scope"].(map[string]interface{})

	auth := "Bearer "
	auth += fmt.Sprint(m["token"])

	url += "/api/alexa/smart_home"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(event))
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err:= ioutil.ReadAll(resp.Body)

	return json.RawMessage(body), err
}

func main() {
        lambda.Start(HandleRequest)
}
