package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// OpenAI API structs compatíveis com o llama.cpp local
type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type Message struct {
	Role       string      `json:"role"`
	Content    string      `json:"content"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	Name       string      `json:"name,omitempty"`
	ToolCalls  interface{} `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float32   `json:"temperature"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role      string     `json:"role"`
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

// Estrutura do arquivo de teste
type Scenario struct {
	ID                     string            `json:"id"`
	Description            string            `json:"description"`
	Tools                  []Tool            `json:"tools"`
	Messages               []Message         `json:"messages"`
	ExpectedToolName       string            `json:"expected_tool_name"`
	ExpectedArgsValidation map[string]string `json:"expected_args_validation"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run tool-bench.go <model-alias>")
		os.Exit(1)
	}
	modelAlias := os.Args[1]
	
	// Carrega os cenários do disco
	data, err := os.ReadFile("bench/testdata/scenarios.json")
	if err != nil {
		fmt.Printf("Erro ao ler scenarios.json: %v\n", err)
		os.Exit(1)
	}

	var scenarios []Scenario
	if err := json.Unmarshal(data, &scenarios); err != nil {
		fmt.Printf("Erro no parse do JSON de teste: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== INICIANDO TOOL BENCH PARA: %s ===\n", modelAlias)
	fmt.Printf("Total de cenários carregados: %d\n\n", len(scenarios))

	// Executa e avalia cada cenário
	for _, sc := range scenarios {
		fmt.Printf("▶ Testando: %s\n", sc.ID)
		
		reqPayload := ChatRequest{
			Model:       modelAlias,
			Messages:    sc.Messages,
			Tools:       sc.Tools,
			Temperature: 0.0,
		}
		
		reqBytes, _ := json.MarshalIndent(reqPayload, "", "  ")
		
		// Faz requisição local pro llama.cpp (Porta 1234)
		resp, err := http.Post("http://localhost:1234/v1/chat/completions", "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			fmt.Printf("❌ Erro HTTP: %v. O servidor está rodando?\n\n", err)
			continue
		}
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		var chatResp ChatResponse
		if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
			fmt.Printf("❌ Erro no parse da resposta do servidor: %v\n", err)
			continue
		}
		
		if len(chatResp.Choices) == 0 {
			fmt.Println("❌ Modelo não retornou nenhuma escolha.")
			fmt.Println()
			continue
		}
		
		msg := chatResp.Choices[0].Message
		
		// Verificações
		if sc.ExpectedToolName == "" {
			// Não esperávamos que chamasse ferramenta
			if len(msg.ToolCalls) == 0 {
				fmt.Println("✅ SUCESSO: Modelo não chamou nenhuma ferramenta (comportamento correto).")
				if msg.Content != "" {
					fmt.Printf("   Resposta Texto: %s\n", msg.Content)
				}
			} else {
				fmt.Printf("❌ FALHA: Modelo alucinou a ferramenta '%s' ao invés de responder texto puro.\n", msg.ToolCalls[0].Function.Name)
			}
		} else {
			// Esperávamos que chamasse ferramenta específica
			if len(msg.ToolCalls) == 0 {
				fmt.Println("❌ FALHA: Modelo respondeu texto ao invés de chamar a ferramenta.")
				fmt.Printf("   Resposta Texto: %s\n", msg.Content)
			} else {
				call := msg.ToolCalls[0]
				if call.Function.Name == sc.ExpectedToolName {
					fmt.Printf("✅ SUCESSO: Chamou a ferramenta certa: %s\n", call.Function.Name)
					fmt.Printf("   Argumentos brutos gerados: %s\n", call.Function.Arguments)
					// TO-DO: Validar os campos do JSON gerado comparando com sc.ExpectedArgsValidation
				} else {
					fmt.Printf("❌ FALHA: Era esperada a ferramenta '%s', mas o modelo chamou '%s'.\n", sc.ExpectedToolName, call.Function.Name)
				}
			}
		}
		fmt.Println("--------------------------------------------------")
	}
}
