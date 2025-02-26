/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chat

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func NewOllamaChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "llama2",                 // 模型名称
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

func NewOpenAIChatModel(ctx context.Context) model.ChatModel {
	key := os.Getenv("OPENAI_API_KEY")
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-4o", // 使用的模型版本
		APIKey: key,      // OpenAI API 密钥
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}

func defaultArkChatModelConfig(ctx context.Context) (*ark.ChatModelConfig, error) {
	config := &ark.ChatModelConfig{
		Model:  os.Getenv("ARK_CHAT_MODEL"),
		APIKey: os.Getenv("ARK_API_KEY"),
	}
	return config, nil
}

func NewArkChatModel(ctx context.Context, config *ark.ChatModelConfig) (cm model.ChatModel) {
	var err error
	if config == nil {
		config, err = defaultArkChatModelConfig(ctx)
		if err != nil {
			log.Fatalf("get default ark config failed, err=%v", err)
		}
	}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		log.Fatalf("create ark chat model failed, err=%v", err)
	}
	return cm
}
