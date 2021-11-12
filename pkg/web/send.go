package web2talk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type AnswerOpenAI struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

type QuestionOpenAI struct {
	Prompt           string   `json:"prompt"`
	MaxTokens        int      `json:"max_tokens"`
	Temperature      float64  `json:"temperature"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	PresencePenalty  float64  `json:"presence_penalty"`
	TopP             float64  `json:"top_p"`
	Stop             []string `json:"stop"`
}

func SendWithArgs(msg string) {

	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		fmt.Println("invalid token, you need export GCP_TOKEN")
		os.Exit(1)
	}

	var question QuestionOpenAI
	question.Prompt = fmt.Sprintf("%v:\n\n1.", msg)
	question.MaxTokens = 64
	question.Temperature = 0.8
	question.FrequencyPenalty = 0.0
	question.PresencePenalty = 0.0
	question.TopP = 1.0
	question.Stop = []string{"\n\n"}

	url := "https://api.openai.com/v1/engines/davinci-instruct-beta/completions"

	b, err := json.Marshal(question)

	body := io.Reader(bytes.NewReader(b))

	req, err := http.NewRequest("POST", url, body)

	// Header with Authorization
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	r := regexp.MustCompile(`20([0-9])`)
	if r.Match([]byte(string(resp.StatusCode))) {
		// if found a problem show the status code
		fmt.Println("statusCode:", resp.StatusCode)
	}

	// read the body with answer
	var answer AnswerOpenAI
	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &answer)
	if err != nil {
		fmt.Println(err)
	}
	// answer.Choices have words of answer
	for _, v := range answer.Choices {
		fmt.Println("<-", v.Text)
	}
}
