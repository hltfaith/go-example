package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
)

const (
	DeepSeekKey  = "xxxxxxxxxxx"
	DeepSeekAPI  = "https://api.deepseek.com/chat/completions"
	DeepSeekMode = "deepseek-chat"
	DoubaoMode   = "doubao-1-5-pro-32k-250115"
	DoubaoKey    = "xxxxxxxxxxx"
	DoubaoAPI    = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
)

type Chat struct {
	Model    string         `json:"model"`
	Messages []*ChatMessage `json:"messages"`
	Stream   bool           `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Post(url string, token string, data []byte) (string, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(60) * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("post request error: %v", err)
		return "", err
	}

	auth := "Bearer " + token
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("post client error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var retData []byte
	retData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("post read body error: %v", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("post error status code: %v", resp.StatusCode)
		return "", fmt.Errorf("err code is :%d", resp.StatusCode)
	}

	ej, _ := simplejson.NewJson(retData)
	res, err := ej.Get("choices").GetIndex(0).Get("message").Get("content").String()
	if err != nil {
		log.Printf("simplejson get error: %v", err)
		return "", err
	}

	return res, nil
}

func main() {
	doubaoMsgList := make([]*ChatMessage, 0)
	doubaoMsgList = append(doubaoMsgList, &ChatMessage{
		Role:    "system",
		Content: "You are a helpful assistant.",
	}, &ChatMessage{
		Role:    "user",
		Content: "现在开始辩论，辩论内容是：你的观点是在北京上班很舒适，，你是反对对方观点的角色，不能认同对方所说的！我的观点在北京上班不舒适，每次辩论不能超过50字。你先开始",
	})
	doubao := &Chat{
		Model:    DoubaoMode,
		Messages: doubaoMsgList,
	}

	deepseekMsgList := make([]*ChatMessage, 0)
	deepseekMsgList = append(deepseekMsgList, &ChatMessage{
		Role:    "system",
		Content: "You are a helpful assistant.",
	}, &ChatMessage{
		Role:    "user",
		Content: "现在开始辩论，辩论内容是：你的观点是认为在北京上班不舒适，必须说在北京上班不舒适的语句，你是反对对方观点的角色！每次回答不能重复，每次辩论不能超过50字。我先开始",
	})
	deepseek := &Chat{
		Model:    DeepSeekMode,
		Messages: deepseekMsgList,
	}

	fmt.Println("现在开始辩论，辩论内容是：豆包大模型的观点是在北京上班很舒适，DeepSeek大模型的观点在北京上班不舒适，每次辩论不能超过50字。豆包大模型先开始\n")

	for {
		body, err := json.Marshal(doubao)
		if err != nil {
			log.Printf("doubao marshal err: %v", err)
			return
		}
		doubaoMsg, err := Post(DoubaoAPI, DoubaoKey, body)
		if err != nil {
			return
		}
		doubaoMsgList = append(doubaoMsgList, &ChatMessage{
			Role:    "assistant",
			Content: doubaoMsg,
		})
		fmt.Printf("[%s] 豆包大模型\n回答：%s\n\n", time.Now().Format("2006-01-02 15:04:05"), doubaoMsg)

		// deepseek
		deepseekMsgList = append(deepseekMsgList, &ChatMessage{
			Role:    "user",
			Content: doubaoMsg,
		})
		body, err = json.Marshal(deepseek)
		if err != nil {
			log.Printf("deepseek marshal err: %v", err)
			return
		}
		deepseekMsg, err := Post(DeepSeekAPI, DeepSeekKey, body)
		if err != nil {
			log.Printf("deepseek post err: %v", err)
			return
		}
		deepseekMsgList = append(deepseekMsgList, &ChatMessage{
			Role:    "assistant",
			Content: deepseekMsg,
		})

		doubaoMsgList = append(doubaoMsgList, &ChatMessage{
			Role:    "user",
			Content: deepseekMsg,
		})

		doubao.Messages = doubaoMsgList
		deepseek.Messages = deepseekMsgList

		fmt.Printf("\x1b[34m[%s] DeepSeek大模型\n回答：%s\x1b[0m\n\n", time.Now().Format("2006-01-02 15:04:05"), deepseekMsg)
	}

}
