package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-ini/ini"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type ChatData struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Max_tokens  int     `json:"max_tokens"`
	Top_p       float64 `json:"top_p"`
	N           float64 `json:"n"`
	Suffix      string  `json:"suffix"`
	Temperature float64 `json:"temperature"`
}

type ChatCompletionResp struct {
	Id      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []CompletionChoices `json:"choices"`
	Usage   Usage               `json:"usage"`
}

type CompletionChoices struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     string `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

const completionUrl = "https://api.openai.com/v1/completions"

const DefaultConfigFile = "~/conf/chatgpt/init.conf"

func ChatCmd() *cobra.Command {
	var content string
	var config string
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Create completion",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(content) <= 0 {
				return errors.New("chat message can not be empty")
			}
			if len(config) == 0 {
				config = DefaultConfigFile
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			createChatCompletion(config, content)
		},
	}
	cmd.Flags().StringVarP(&content, "message", "m", "", "Chat with ChatGPT message")
	cmd.Flags().StringVarP(&config, "config", "c", "", "ChatGPT config")
	cmd.MarkFlagRequired("message")
	return cmd

}

func createChatCompletion(configPath string, content string) {
	config := loadConfig(configPath)
	openApiKey, exist := config["openai_api_key"].(string)
	if !exist {
		panic("openai api key can not be null")
	}

	chatData := &ChatData{
		Max_tokens:  config["max_tokens"].(int),
		Model:       config["model"].(string),
		Prompt:      content,
		Top_p:       config["top_p"].(float64),
		N:           config["n"].(float64),
		Suffix:      config["suffix"].(string),
		Temperature: config["temperature"].(float64),
	}
	proxy, exist := config["proxy"].(string)
	if !exist {
		proxy = loadProxy()
	}
	var client http.Client
	if len(proxy) > 0 {
		parse, err := url.Parse(proxy)
		if err != nil {
			panic(err)
		}
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(parse),
			},
		}
	}

	marshal, err2 := json.Marshal(chatData)
	if err2 != nil {
		panic(err2)
	}
	s := string(marshal)
	//	var body = strings.NewReader(s)
	request, err := http.NewRequest("POST", completionUrl, bytes.NewBufferString(s))
	request.Header.Add("Authorization", "Bearer "+openApiKey)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	all, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		panic(err2)
	}

	chatResp := ChatCompletionResp{}

	_ = json.Unmarshal(all, &chatResp)

	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"result"})
	writer.AppendRow(table.Row{chatResp.Choices[0].Text})
	writer.Render()
}

func loadProxy() string {
	return ""
}

func loadConfig(configPath string) map[string]any {
	load, err := ini.Load(configPath)
	if err != nil {
		panic(err)
	}

	gptSection, err := load.GetSection("gpt")
	if err != nil {
		panic(err)
	}

	return map[string]any{
		"n":              gptSection.Key("n").MustFloat64(1),
		"model":          gptSection.Key("model").MustString("text-davinci-003"),
		"top_p":          gptSection.Key("top_p").MustFloat64(1),
		"proxy":          gptSection.Key("proxy").Value(),
		"suffix":         gptSection.Key("suffix").MustString(""),
		"max_tokens":     gptSection.Key("max_tokens").MustInt(3000),
		"temperature":    gptSection.Key("temperature").MustFloat64(0.2),
		"openai_api_key": gptSection.Key("openai_api_key").Value(),
	}
}
