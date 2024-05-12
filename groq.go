package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const groqURL = "https://api.groq.com/openai/v1/chat/completions"
const bodyStart = `{"messages": [{"role": "user", "content": "`
const bodyEnd = `"}], "model": "`
const bodyEndReal = `"}`

func AskGroq(question string) string {

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + *groqKey

	var bodyReader = strings.NewReader(bodyStart + question + bodyEnd + *groqModel + bodyEndReal)

	// Create a new request using http
	req, err := http.NewRequest("POST", groqURL, bodyReader)
	if err != nil {
		panic(err)
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	var p GroqReponse
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", p)

	return p.Choices[0].Message.Content
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
