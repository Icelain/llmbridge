package llmbridge

import (
	"context"
	ollama "github.com/ollama/ollama/api"
	"log"
)

type LLMClient interface {
	StreamResponse(ctx context.Context, query string) <-chan string
}

type DefaultLLMClient struct {
	ollamaclient *ollama.Client
	model        string
}

func (dl *DefaultLLMClient) StreamResponse(ctx context.Context, query string) <-chan string{

	to := make(chan string)
	request := &ollama.GenerateRequest{

		Model:  dl.model,
		Prompt: query,
	}

	defer func(){dl.ollamaclient.Generate(ctx, request, func(generatedResponse ollama.GenerateResponse) error {

		to <- generatedResponse.Response
		return nil

	})

	close(to)}()

	return to

}

func NewDefaultLLMClient(ctx context.Context, model string) *DefaultLLMClient {

	client, err := ollama.ClientFromEnvironment()
	if err != nil {

		log.Fatal(err)

	}

	return &DefaultLLMClient{ollamaclient: client, model: model}

}
