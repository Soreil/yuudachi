package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

const groqURL = "https://api.groq.com/openai/v1/chat/completions"

var CurrentTemperature float64 = 0.6
var CurrentReasoningFormat reasoningFormats = parsedReasoningFormat

func AskGroqSystem(prompt string, context []Message) []Message {
	if context == nil {
		context = []Message{}
	}

	context = append(context, Message{
		Role:    "system",
		Content: prompt,
	})

	return context
}

func AskGroq(question string, context []Message) (string, []Message, error) {

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + *groqKey

	if context == nil {
		context = []Message{}
	}

	context = append(context, Message{
		Role:    "user",
		Content: question,
	},
	)

	var groqRequest = GroqRequest{
		context,
		*groqModel,
		CurrentTemperature,
		CurrentReasoningFormat,
	}

	data, err := json.Marshal(groqRequest)

	if err != nil {
		log.Println(err)
		return "", context, err
	}

	var reader = bytes.NewReader(data)

	// Create a new request using http
	req, err := http.NewRequest("POST", groqURL, reader)
	if err != nil {
		log.Println(err)
		return "", context, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err2 := fmt.Errorf("Error on response.\n[ERROR] -: %v", err)

		return "", context, err2
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", context, errors.New(resp.Status)
	}

	var p GroqReponse
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", p)

	if len(p.Choices) == 0 {
		return "", context, processError(resp)
	}
	result := p.Choices[len(p.Choices)-1].Message
	context = append(context, result)
	return result.Content, context, nil
}

func processError(resp *http.Response) error {
	data, err := httputil.DumpResponse(resp, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q", data)

	return errors.New("no responses message in body")
}

type reasoningFormats string

const (
	rawReasoningFormat    reasoningFormats = "raw"
	parsedReasoningFormat reasoningFormats = "parsed"
	hiddenReasoningFormat reasoningFormats = "hidden"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Messages        []Message        `json:"messages"`
	Model           string           `json:"model"`
	Temperature     float64          `json:"temperature"`
	ReasoningFormat reasoningFormats `json:"reasoning_format"`
}

type GroqReponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
	XGroq             struct {
		ID string `json:"id"`
	} `json:"x_groq"`
}
