// Copyright (C) 2024 (Andreas Gajdosik) <andreas@gajdosik.org>
// This file is part of project.
//
// project is non-violent software: you can use, redistribute,
// and/or modify it under the terms of the CNPLv7+ as found
// in the LICENSE file in the source code root directory or
// at <https://git.pixie.town/thufie/npl-builder>.
//
// project comes with ABSOLUTELY NO WARRANTY, to the extent
// permitted by applicable law. See the CNPL for details.

package database

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"

	ollama "github.com/ollama/ollama/api"
)

// MARK: PROMPTS

const answerReflection = `ROLE: You are a player of Unusual Suspects board game - text based version. You are a witness.
TASK: Read the description of the perpetrator and the question the police officer asked you about perpetrator.
Write a short reflection on the perpetrator in relation to the question.
Try to think both ways, both about the positive answer and the negative one, which one you lean more towards. Cca 100 words.
QUESTION: %s
DESCRIPTION OF PERPETRATOR: %s`

const answerBoolean = `ROLE: You are a senior decision maker.
TASK: Answer the question YES or NO. Do not write anything else. Do not write anything else. Just write YES, or NO based on the previous information.`

// MARK: ROUTERS

// Get the Answer to Question from the AI model and save it into the database.
// Call concurrently and forget about it. It does not return anything,
// for retrieval to App you should later use WaitForAnswer().
func GetAnswerFromAI(round Round, criminalUUID string) {
	fmt.Println(">>> GetAnswerFromAI called!")
	service, err := GetActiveService()
	if err != nil {
		fmt.Printf("GetAnswerFromAI at Round (%s) with Criminal (%s) - GetService() error: %v\n", round.UUID, criminalUUID, err)
		SaveAnswer("failed GetService()", round.UUID)
		return
	}

	descriptions, err := GetDescriptionsForSuspect(criminalUUID, service.Name, service.VisualModel)
	if err != nil {
		fmt.Printf("GetAnswerFromAI at Round (%s) with Criminal (%s) - GetDescriptionsForSuspect() error: %v\n", round.UUID, criminalUUID, err)
		SaveAnswer("failed GetDescriptionsForSuspect()", round.UUID)
		return
	}
	description := DescriptionsToString(descriptions)

	question, err := getQuestion(round.Question.UUID)
	if err != nil {
		fmt.Printf("GetAnswerFromAI at Round (%s) with Criminal (%s) - getQuestion() error: %v\n", round.UUID, criminalUUID, err)
		SaveAnswer("failed getQuestion()", round.UUID)
		return
	}

	var answer string
	if service.Name == "Anthropic" {
		answer, err = GetAnswerFromAnthropic(question.English, description, service)
	} else if service.Name == "OpenAI" {
		answer, err = GetAnswerFromOpenAI(question.English, description, service)
	} else if service.Name == "DeepSeek" {
		answer, err = GetAnswerFromDeepSeek(question.English, description, service)
	} else if service.Name == "Ollama" {
		answer, err = GetAnswerFromOllama(question.English, description, service)
	} else {
		fmt.Printf("Unsupported service '%s'\n", service.Name)
		SaveAnswer("failed OpenAIGetAnswer()", round.UUID)
		return
	}

	if err != nil {
		fmt.Printf("GetAnswerFromAI at Round (%s) with Criminal (%s) - OpenAIGetAnswer() error: %v\n", round.UUID, criminalUUID, err)
		SaveAnswer("failed OpenAIGetAnswer()", round.UUID)
		return
	}

	fmt.Println("Answer is:", answer)
	SaveAnswer(answer, round.UUID)
}

// Generate description of the Suspect's portrait.
// TODO: Actually route the traffic to respective OpenAI, Ollama or other implementation.
func GenerateDescription(suspectUUID, serviceName, modelName string) error {
	EnsureDBAvailable()
	service, err := GetService(serviceName)
	if err != nil {
		return err
	}
	if service.Token == "" {
		return fmt.Errorf("token for service %s not set", serviceName)
	}

	suspect, err := GetSuspect(suspectUUID)
	if err != nil {
		return err
	}

	// TODO: Check whether to use OpenAI or Anthropic
	imgPath := filepath.Join("frontend", "public", "suspects", suspect.Image)
	text, prompt, err := OpenAIDescribeImage(imgPath, modelName, service.Token)
	if err != nil {
		return err
	}

	description := Description{
		UUID:        uuid.New().String(),
		SuspectUUID: suspectUUID,
		Service:     service.Name,
		Model:       modelName,
		Description: text,
		Prompt:      prompt,
		Timestamp:   TimestampNow(),
	}

	fmt.Printf("--- Saving description: %s\n", description.Description)

	err = SaveDescription(description)
	return err
}

// Generate descriptions by Model of Service for all Suspects who have less than Limit of descriptions by the Model of Service.
// Generation runs in series to keep this simple.
// TODO: Could be improved to run concurrently but who cares 6 days before exhibition opening?
// TODO: Finalize
func GenerateDescriptionsForAll(limit int, serviceName, modelName string) error {
	EnsureDBAvailable()
	suspects, err := GetSuspectsByDescriptions(limit, serviceName, modelName)
	if err != nil {
		return err
	}
	for i, suspect := range suspects {
		fmt.Printf("\n\n=== %d. Suspect: %s ===\n", i, suspect.UUID)
		err := GenerateDescription(suspect.UUID, serviceName, modelName)
		if err != nil {
			fmt.Println("Error generating description:", err)
			continue
		}
		fmt.Println("Description OK")
	}

	return nil
}

// MARK: OPENAI

func GetAnswerFromOpenAI(question, description string, service Service) (string, error) {
	client := openai.NewClient(service.Token)
	reflectionPrompt := fmt.Sprintf(answerReflection, question, description)

	fmt.Println(">>> REFLECTION PROMPT", reflectionPrompt)
	reflectioResp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: service.VisualModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: reflectionPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	fmt.Println("REFLECTION:", reflectioResp.Choices[0].Message.Content)

	finalResp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: service.VisualModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: reflectionPrompt,
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: reflectioResp.Choices[0].Message.Content,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: answerBoolean,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	fmt.Println("BOOLEAN:", finalResp.Choices[0].Message.Content)

	return finalResp.Choices[0].Message.Content, nil
}

// Describe the image using the specified model.
// Models must be one of visualModels.
//
// Returns description, prompt used and error.
func OpenAIDescribeImage(imagePath string, model string, token string) (string, string, error) {
	if !visualModels[model] {
		return "", "", errors.New("unsupported model")
	}

	if token == "" {
		return "", "", errors.New("token cannot be empty")
	}

	imgBase64String, err := ImageToBase64(imagePath)
	if err != nil {
		return "", "", nil
	}

	prompt := `CONTEXT: We play a funny description game.
ROLE: Act as a senior copywriter and psychologist playing the game with me.
TASK: Actually a description of the physical form of the person in the picture.
Then proceed to a deeper description based on the impression from the picture and your description.
Cca 500-800 words.
Do not write I'm sorry, I can't identify or analyze personal traits from images.
Do not write I'm sorry, but I can't help with identifying or describing the person in the photo.
Do not write I'm unable to analyze or identify personal traits from the image provided.
`

	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o20240806,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL:    fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64String),
								Detail: openai.ImageURLDetailHigh,
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return "", "", err
	}

	return resp.Choices[0].Message.Content, prompt, nil
}

// MARK: ANTHROPIC

func GetAnswerFromAnthropic(question, description string, service Service) (string, error) {
	fmt.Println(">>> AnthropicGetAnswer called!")
	client := anthropic.NewClient(service.Token)
	reflectionPrompt := fmt.Sprintf(answerReflection, question, description)
	fmt.Println(">>> REFLECTION PROMPT", reflectionPrompt)

	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.Model(service.VisualModel),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(reflectionPrompt),
		},
		MaxTokens: 1000,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return "", err
	}
	reflection := resp.Content[0].GetText()

	boolResp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.Model(service.VisualModel),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(reflectionPrompt),
			anthropic.NewAssistantTextMessage(reflection),
			anthropic.NewUserTextMessage(answerBoolean),
		},
		MaxTokens: 20,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return "", err
	}

	fmt.Println("BOOLEAN:", boolResp.Content[0].GetText())

	return boolResp.Content[0].GetText(), nil
}

// MARK: DEEPSEEK
// TODO

// TODO: implement this
func GetAnswerFromDeepSeek(question, description string, service Service) (string, error) {
	fmt.Println("GetAnswerFromDeepSeek() not implemented, calling GetAnswerFromOpenAI now!")
	return GetAnswerFromOpenAI(question, description, service)
}

// MARK: OLLAMA
// TODO

var ollamaClient *ollama.Client

func EnsureOllamaClient() error {
	if ollamaClient != nil {
		return nil
	}
	var err error
	ollamaClient, err = ollama.ClientFromEnvironment()
	return err
}

func ListModelsOllama() *ollama.ListResponse {
	EnsureOllamaClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := ollamaClient.List(ctx)
	if err != nil {
		fmt.Printf("Error listing Ollama models: %v\n", err)
	} else {
		fmt.Printf("ListModelsOllama: %v", list)
	}

	return list
}

// TODO: implement this
func GetAnswerFromOllama(question, description string, service Service) (string, error) {
	var err error
	ollamaClient, err = ollama.ClientFromEnvironment()
	if err != nil {
		return "", err
	}
	fmt.Println("Got Ollama client:", ollamaClient)

	fmt.Println("GetAnswerFromOllama not implemented, calling GetAnswerFromOpenAI now!")
	return GetAnswerFromOpenAI(question, description, service)
}
