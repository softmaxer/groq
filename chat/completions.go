package chat

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ChatCompletionError struct {
	StatusCode int
}

func (c *ChatCompletionError) Error() string {
	return fmt.Sprintf("Error from completion response API. Status code: %d\n", c.StatusCode)
}

func GetChatCompletion(
	chatRequest *ChatRequest,
	apiKey string,
) (ChatResponse, error) {
	var bearer string = fmt.Sprintf("Bearer %s", apiKey)
	encodedRequest, err := encode(chatRequest)
	if err != nil {
		return ChatResponse{}, err
	}
	bodyReader := bytes.NewReader(encodedRequest)
	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bodyReader,
	)
	if err != nil {
		log.Println(err.Error())
		return ChatResponse{}, err
	}
	req.Header.Add("Authorization", bearer)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ChatResponse{}, err
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Got status code: %d, exiting with empty response", response.StatusCode)
		return ChatResponse{}, &ChatCompletionError{StatusCode: response.StatusCode}
	}
	defer response.Body.Close()
	chatResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return ChatResponse{}, err
	}
	var objMap ChatResponse
	err = decode(chatResponse, &objMap)
	if err != nil {
		return ChatResponse{}, err
	}
	return objMap, nil
}

func CreateCompletionRequest(modelName, system, userRequest string) ChatRequest {
	systemMessage := Message{Role: "system", Content: system}
	userMessage := Message{Role: "user", Content: userRequest}
	return ChatRequest{Messages: []Message{systemMessage, userMessage}, Model: modelName}
}
