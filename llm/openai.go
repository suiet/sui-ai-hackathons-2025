package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"suiet_server/utils/config"
)

var LLM *openai.Client

func init() {
	if config.GetServerType() != "api" {
		return
	}
	LLM = openai.NewClient(
		option.WithBaseURL("https://api.model.box/v1/"),  // defaults to https://api.openai.com
		option.WithAPIKey(os.Getenv("MODELBOX_API_KEY")), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
}

func ValidateNFTImage(title, description, imageUrl string) (isSafe bool, invalidateCategories []string, err error) {

	t := template.Must(template.New("validate").Parse(validateNFTImage))
	data := struct {
		Title       string
		Description string
		ImageURL    string
	}{
		Title:       title,
		Description: description,
		ImageURL:    imageUrl,
	}

	var sb strings.Builder

	err = t.Execute(&sb, data)
	if err != nil {
		return false, nil, err
	}

	imageBase64, err := fetchImageAsBase64(imageUrl)
	if err != nil {
		return false, nil, err
	}

	chatCompletion, err := LLM.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sb.String()),
			openai.UserMessageParts(openai.ImagePart(imageBase64)),
		}),
		//Model: openai.F("qwen/qwen2-vl-72b-instruct"),
		Model:       openai.F("openai/gpt-4o-mini"),
		Temperature: openai.Float(0),
	})
	if err != nil {
		return false, nil, err
	}

	result := chatCompletion.Choices[0].Message.Content
	fmt.Println(result)
	resultLines := strings.Split(result, "\n")
	if len(resultLines) == 0 {
		return false, nil, fmt.Errorf("no result")
	}

	if strings.TrimSpace(resultLines[0]) != "safe" && strings.TrimSpace(resultLines[0]) != "unsafe" {

		return false, nil, fmt.Errorf("error parse result")
	}
	if strings.TrimSpace(resultLines[0]) == "safe" {
		isSafe = true
		return isSafe, nil, nil
	}

	// must be unsafe now
	//if strings.TrimSpace(resultLines[0]) == "unsafe" {
	isSafe = false

	if len(resultLines) > 1 {
		invalidateCategories = strings.Split(resultLines[1], ",")
	}
	return isSafe, invalidateCategories, nil
	//}

}

func fetchImageAsBase64(url string) (string, error) {
	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: HTTP %d", resp.StatusCode)
	}

	// Get Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Read image content
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image content: %v", err)
	}

	// Encode image content to Base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Generate HTML-compatible src format
	base64Src := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
	return base64Src, nil
}
