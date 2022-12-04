package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"os"
	"tailscale.com/tsnet"
	"time"
)

func HandleRequest(ctx context.Context, event json.RawMessage) (json.RawMessage, error) {
	url := os.Getenv("BASE_URL")
	token_env, token_present := os.LookupEnv("TOKEN")

	auth := "Bearer "
	// use token from environment if present
	if token_present {
		auth += token_env
	} else {
		var f interface{}
		_ = json.Unmarshal(event, &f)
		m := f.(map[string]interface{})["directive"].(map[string]interface{})["endpoint"].(map[string]interface{})["scope"].(map[string]interface{})
		auth += fmt.Sprint(m["token"])
	}

	url += "/api/alexa/smart_home"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(event))
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	s := &tsnet.Server{
		Dir:       "/tmp",
		Hostname:  "hass-alexa",
		Ephemeral: true,
	}

	defer s.Close()

	_, _ = s.LocalClient()
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: s.Dial, // use the tailscale dialer
		},
	}
	time.Sleep(1500 * time.Millisecond) // wait for tailnet connection to come up

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return json.RawMessage(body), err
}

func main() {
	lambda.Start(HandleRequest)
}
