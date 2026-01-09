package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"smartcommunity/internal/config"
)

type AIService struct{}

// 阿里云 DashScope 请求结构
type DashScopeRequest struct {
	Model    string             `json:"model"`
	Messages []DashScopeMessage `json:"messages"`
}

type DashScopeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 阿里云 DashScope 响应结构 (简化)
type DashScopeResponse struct {
	Choices []struct {
		Message DashScopeMessage `json:"message"`
	} `json:"choices"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"` // 错误信息
}

// Chat 调用 AI 对话
func (s *AIService) Chat(userContent string) (string, error) {
	apiKey := config.Conf.AI.APIKey
	baseURL := config.Conf.AI.BaseURL
	modelName := config.Conf.AI.Model

	if apiKey == "" || baseURL == "" {
		return "", errors.New("AI 服务未配置")
	}

	// 构建请求体
	reqBody := DashScopeRequest{
		Model: modelName,
		Messages: []DashScopeMessage{
			{Role: "system", Content: "你是'智慧社区'APP的专属智能助手。你的职责是帮助用户解决社区生活相关的问题，例如：缴纳物业费、报事报修、查看社区公告、访客预约、车位管理、社区商城购物等。请拒绝回答与社区管理、生活服务、商城购物无关的问题，并礼貌地表明身份。"},
			{Role: "user", Content: userContent},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var aiResp DashScopeResponse
	if err := json.Unmarshal(bodyText, &aiResp); err != nil {
		return "", err
	}

	if aiResp.Code != "" {
		// API 返回了错误码 (DashScope 错误时会有 Code 字段)
		return "", errors.New("AI API Error: " + aiResp.Message)
	}

	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "", errors.New("No response from AI")
}
