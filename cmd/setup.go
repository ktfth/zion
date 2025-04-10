package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"zion/config"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configura o ambiente do Zion",
	Long:  `Configura o diretório home do Zion e instala o plugin HelloWorld de exemplo.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Carregar a configuração
		cfg := config.LoadConfig()

		fmt.Printf("Configurando o ambiente Zion em: %s\n", cfg.HomeDir)

		// Criar o diretório de plugins se não existir
		if err := os.MkdirAll(cfg.PluginsDir, 0755); err != nil {
			fmt.Printf("Erro ao criar diretório de plugins: %v\n", err)
			return
		}

		// Criar arquivo README.md no diretório home
		readmePath := filepath.Join(cfg.HomeDir, "README.md")
		readmeContent := `# Zion Home Directory

Este é o diretório home do Zion, onde são armazenadas configurações e plugins.

## Estrutura

- plugins/ - Diretório onde os plugins são armazenados (apenas para sistemas Unix)
- config.yaml - Arquivo de configuração (opcional)

## Plugins no Windows

No Windows, o Go não suporta plugins dinâmicos (buildmode=plugin). Portanto, os plugins no Windows são implementados estaticamente no código fonte do Zion.

Para adicionar um novo plugin no Windows:

1. Crie um novo arquivo .go no diretório plugins/ do código fonte do Zion
2. Implemente a interface Plugin
3. Registre o plugin usando RegisterPlugin() na função init()
4. Recompile o Zion

## Plugins em sistemas Unix (Linux/macOS)

Em sistemas Unix, os plugins podem ser compilados como arquivos .so e colocados no diretório plugins/.
Exemplo de como compilar um plugin:

` + "```bash" + `
go build -buildmode=plugin -o hello_world.so hello_world.go
` + "```" + `

## Plugin HelloWorld

O plugin HelloWorld já está incluído estaticamente no Zion e demonstra como criar plugins.
Ele adiciona instruções extras ao prompt de geração de scaffold e exibe mensagens durante o processo.
`

		if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
			fmt.Printf("Erro ao criar arquivo README.md: %v\n", err)
			return
		}

		// Criar diretório para o plugin HelloWorld de exemplo
		helloWorldDir := filepath.Join(cfg.HomeDir, "plugins", "hello_world")
		if err := os.MkdirAll(helloWorldDir, 0755); err != nil {
			fmt.Printf("Erro ao criar diretório do plugin HelloWorld: %v\n", err)
			return
		}

		// Criar o arquivo hello_world.go
		helloWorldPath := filepath.Join(helloWorldDir, "hello_world.go")
		helloWorldContent := `package main

import (
	"fmt"
	"strings"
)

// ScaffoldHook define os pontos de extensão para plugins
type ScaffoldHook string

const (
	// BeforeGeneration é executado antes da geração do scaffold
	BeforeGeneration ScaffoldHook = "before_generation"
	// AfterGeneration é executado após a geração do scaffold
	AfterGeneration ScaffoldHook = "after_generation"
	// ModifyPrompt permite modificar o prompt antes de enviá-lo para a API
	ModifyPrompt ScaffoldHook = "modify_prompt"
)

// ScaffoldContext contém informações sobre o processo de geração de scaffold
type ScaffoldContext struct {
	ProjectName string
	Language    string
	Description string
	Prompt      string
	Response    string
}

// Plugin define a interface que todo plugin deve implementar
type Plugin interface {
	Name() string
	Execute() error
	GetHooks() map[ScaffoldHook]interface{}
}

// HelloWorldPlugin é um plugin de exemplo que adiciona uma mensagem de boas-vindas
// e modifica o prompt para incluir requisitos adicionais
type HelloWorldPlugin struct{}

// Name retorna o nome do plugin
func (p *HelloWorldPlugin) Name() string {
	return "HelloWorld"
}

// Execute é chamado quando o plugin é executado
func (p *HelloWorldPlugin) Execute() error {
	fmt.Println("HelloWorld plugin: Olá, mundo!")
	return nil
}

// GetHooks retorna os hooks que o plugin implementa
func (p *HelloWorldPlugin) GetHooks() map[ScaffoldHook]interface{} {
	return map[ScaffoldHook]interface{}{
		BeforeGeneration: p.beforeGeneration,
		ModifyPrompt:     p.modifyPrompt,
		AfterGeneration:  p.afterGeneration,
	}
}

// beforeGeneration é executado antes da geração do scaffold
func (p *HelloWorldPlugin) beforeGeneration(ctx *ScaffoldContext) error {
	fmt.Printf("HelloWorld plugin: Iniciando geração para projeto '%s' em %s\n", 
		ctx.ProjectName, ctx.Language)
	return nil
}

// modifyPrompt modifica o prompt para incluir requisitos adicionais
func (p *HelloWorldPlugin) modifyPrompt(prompt string) string {
	// Adiciona requisitos específicos do HelloWorld plugin
	additionalRequirements := "\n\nAdicional do HelloWorld Plugin:\n" +
		"1. Adicione um arquivo hello.md com uma mensagem de boas-vindas\n" +
		"2. Inclua comentários explicativos no código\n"
	
	// Insere os requisitos adicionais antes da linha IMPORTANTE
	if idx := strings.Index(prompt, "IMPORTANTE:"); idx != -1 {
		return prompt[:idx] + additionalRequirements + prompt[idx:]
	}
	
	// Se não encontrar o marcador, apenas adiciona ao final
	return prompt + additionalRequirements
}

// afterGeneration é executado após a geração do scaffold
func (p *HelloWorldPlugin) afterGeneration(ctx *ScaffoldContext) error {
	fmt.Printf("HelloWorld plugin: Geração concluída para projeto '%s'\n", ctx.ProjectName)
	return nil
}

// Plugin é a variável exportada que o Zion irá procurar
var Plugin HelloWorldPlugin
`

		if err := os.WriteFile(helloWorldPath, []byte(helloWorldContent), 0644); err != nil {
			fmt.Printf("Erro ao criar arquivo hello_world.go: %v\n", err)
			return
		}

		// Criar arquivo README.md para o plugin HelloWorld
		helloWorldReadmePath := filepath.Join(helloWorldDir, "README.md")
		helloWorldReadmeContent := `# Plugin HelloWorld

Este é um plugin de exemplo para o Zion que demonstra como criar plugins que podem modificar o processo de geração de scaffold.

## Funcionalidades

- Adiciona uma mensagem de boas-vindas ao gerar um projeto
- Modifica o prompt para incluir requisitos adicionais
- Demonstra como usar os hooks BeforeGeneration, ModifyPrompt e AfterGeneration

## Compilação

Para compilar o plugin, execute:

` + "```bash" + `
go build -buildmode=plugin -o hello_world.so hello_world.go
` + "```" + `

## Instalação

Copie o arquivo hello_world.so para o diretório ~/.zion/plugins/:

` + "```bash" + `
cp hello_world.so ~/.zion/plugins/
` + "```" + `

## Uso

Após a instalação, o plugin será carregado automaticamente quando você executar o comando zion scaffold.
`

		if err := os.WriteFile(helloWorldReadmePath, []byte(helloWorldReadmeContent), 0644); err != nil {
			fmt.Printf("Erro ao criar arquivo README.md do plugin HelloWorld: %v\n", err)
			return
		}

		fmt.Println("Configuração concluída com sucesso!")
		fmt.Println("Diretório home do Zion:", cfg.HomeDir)
		fmt.Println("Diretório de plugins:", cfg.PluginsDir)
		fmt.Println("Plugin HelloWorld de exemplo criado em:", helloWorldDir)
		
		if runtime.GOOS == "windows" {
			fmt.Println("\nNota: No Windows, os plugins são implementados estaticamente.")
			fmt.Println("O plugin HelloWorld já está incluído no código fonte do Zion.")
			fmt.Println("Para criar novos plugins, adicione-os ao diretório plugins/ do código fonte e recompile o Zion.")
		} else {
			fmt.Println("\nPara compilar e instalar o plugin HelloWorld, execute:")
			fmt.Printf("cd %s\n", helloWorldDir)
			fmt.Println("go build -buildmode=plugin -o hello_world.so hello_world.go")
			fmt.Printf("cp hello_world.so %s\n", cfg.PluginsDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
