package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/gtsteffaniak/ai-fake-news/routes"
	"google.golang.org/api/option"
)

var (
	wg       sync.WaitGroup
	articles = map[string]bool{}
	ctx      = context.Background()
	model    *genai.GenerativeModel
)

func main() {
	model = setupLLMClient() // ai client
	devMode := flag.Bool("dev", false, "enable dev mode (hot-reloading and debug logging)")
	flag.Parse()
	opts := &slog.HandlerOptions{
		// Use the ReplaceAttr function on the handler options
		// to be able to replace any single attribute in the log output
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// check that we are handling the time key
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()
			// change the value from a time.Time to a String
			// where the string has the correct time format.
			a.Value = slog.StringValue(t.Format(time.DateTime))
			return a
		},
	}
	if *devMode {
		opts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	slog.Debug("Program was run in dev mode which enables debugging and hotloading.")
	wg.Add(1)
	go func() {
		defer wg.Done()
		routes.SetupWeb(*devMode, *logger)
	}()
	wg.Wait()

	// Call Gemini API with sanitizedText
	//r, err := getLLMResponse(model)
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(string(r))

}

// This function is not includedin the provided code, but represents the logic for calling the Gemini API
func getLLMResponse(model *genai.GenerativeModel) (string, error) {
	fmt.Println("getting model respoonse")
	prompt := `
	output json format array of article objects:
	[
	  {"title": "new story title here", "article": "<p> A report about the topic</p>" "summary":"one or two sentance summary of article","category":"technology"}
	]

	each article has:
	"title" is always news headline, 10 words max. generate based on story below.
	"article" is html content of the generated article. (minimum 5 sentences and two paragraphs.)
	"category" is the topic cateogry the article is in
	"summary" is a one or two sentence summary of what the article is about

	Generate stories:
	1. technology: a dog robot named spot has been infected with rust the programming language
	2. politics: the house of representatives have decided to take a sebatical for the next year.
	3. science: nano cell technology has become cheap and affordable.
	`
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	responseString := ""
	// Show the model's response, which is expected to be text.
	for _, part := range resp.Candidates[0].Content.Parts {
		responseString += fmt.Sprintf("%v", part)
	}
	responseString = strings.Trim(responseString, "`")
	responseString = strings.Replace(responseString, "`", "", -1)
	responseString = strings.Trim(responseString, "json")
	responseString = strings.Trim(responseString, "\n")
	// Return an empty string if no response is found
	return responseString, err
}

func setupLLMClient() *genai.GenerativeModel {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel("gemini-1.5-flash") // change model?
	return model
}
