package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
)

type ChatDataConfig struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	TopP        float64 `json:"top_p"`
	N           float64 `json:"n"`
	Suffix      string  `json:"suffix"`
	Temperature float64 `json:"temperature"`
	OpenApiKey  string  `json:"-"`
	Proxy       string
}

type ChatDataRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	TopP        float64 `json:"top_p"`
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

const DefaultConfigFile = "%s/conf/chatgpt/init.conf"

func ChatCmd() *cobra.Command {
	var content string
	var config string
	var key string
	var model string
	var proxy string
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Create completion",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(content) <= 0 {
				return errors.New("chat message can not be empty")
			}
			if len(config) == 0 {
				current, err := user.Current()
				if err != nil {
					log.Fatal(err)
					return err
				}
				config = fmt.Sprintf(DefaultConfigFile, current.HomeDir)
			}
			_, err := os.Stat(config)
			if os.IsNotExist(err) && len(key) == 0 {
				return errors.New(fmt.Sprintf("You must specify a key or a configuration file. The default configuration file path of %s does not exist ", config))
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			configData := loadConfig(config)
			if configData == nil {
				configData = &ChatDataConfig{}
			}
			if len(key) > 0 {
				configData.OpenApiKey = key
			}
			createChatCompletion(configData, content)
		},
	}
	cmd.Flags().StringVarP(&content, "message", "m", "", "Chat with ChatGPT message")
	cmd.Flags().StringVarP(&config, "config", "c", "", "ChatGPT config")
	cmd.Flags().StringVarP(&key, "key", "k", "", "ChatGPT Open API key. If no configuration file is specified, the configuration file is used. If no configuration file exists, an error is reported")
	cmd.Flags().StringVarP(&model, "model", "d", "", "ChatGPT ID of the Model to use . Default is text-davinci-003")
	cmd.Flags().StringVarP(&proxy, "proxy", "p", "", "Network proxy.")
	cmd.MarkFlagRequired("message")
	return cmd

}

func createChatCompletion(chatDataConfig *ChatDataConfig, content string) {
	if len(chatDataConfig.OpenApiKey) == 0 {
		panic("openai api key can not be null")
	}
	chatData := &ChatDataRequest{
		MaxTokens:   chatDataConfig.MaxTokens,
		Model:       chatDataConfig.Model,
		Prompt:      content,
		TopP:        chatDataConfig.TopP,
		N:           chatDataConfig.N,
		Suffix:      chatDataConfig.Suffix,
		Temperature: chatDataConfig.Temperature,
	}

	proxy := chatDataConfig.Proxy
	if len(proxy) == 0 {
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
	request.Header.Add("Authorization", "Bearer "+chatDataConfig.OpenApiKey)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		log.Fatal(response.StatusCode)
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

func loadConfig(configPath string) *ChatDataConfig {
	load, err := ini.Load(configPath)
	if os.IsNotExist(err) {
		return nil
	}
	gptSection, err := load.GetSection("gpt")
	if err != nil {
		panic(err)
	}
	return &ChatDataConfig{
		N:           gptSection.Key("n").MustFloat64(1),
		Model:       gptSection.Key("model").MustString("text-davinci-003"),
		TopP:        gptSection.Key("top_p").MustFloat64(1),
		Proxy:       gptSection.Key("proxy").MustString(""),
		Suffix:      gptSection.Key("suffix").MustString(""),
		MaxTokens:   gptSection.Key("max_tokens").MustInt(3000),
		Temperature: gptSection.Key("temperature").MustFloat64(0.2),
		OpenApiKey:  gptSection.Key("openai_api_key").MustString(""),
	}
}
