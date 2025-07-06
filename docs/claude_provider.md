# 🧠 Claude Provider - Documentação Técnica

## Visão Geral

O Zion CLI oferece suporte completo ao Claude (Anthropic) através de duas implementações:

1. **Claude via OpenRouter** - Usando o provider OpenRouter com modelos Claude
2. **Claude Provider Dedicado** - Provider específico para Claude

## Arquitetura

### Estrutura de Arquivos

```
ai/providers/
├── claude.go          # Provider Claude dedicado
├── openrouter.go      # Provider OpenRouter (também suporta Claude)
├── manager.go         # Gerenciador de providers
└── provider.go        # Interface base
```

### Registro de Providers

```go
// ai/providers/manager.go
func NewProviderManager() *ProviderManager {
    manager := &ProviderManager{
        factories: make(map[string]Factory),
        providers: make(map[string]Provider),
    }
    
    // Registro de providers
    manager.RegisterFactory("gemini", NewGeminiProvider)
    manager.RegisterFactory("gpt", NewGPTProvider)
    manager.RegisterFactory("openrouter", NewOpenRouterProvider)
    manager.RegisterFactory("claude", NewClaudeProvider)  // ✅ Novo
    
    return manager
}
```

## Implementação do Claude Provider

### Estrutura Principal

```go
// ai/providers/claude.go
type ClaudeProvider struct {
    apiKey      string
    model       string
    maxTokens   int
    temperature float64
    baseURL     string
}
```

### Função de Criação

```go
func NewClaudeProvider(config map[string]string) (Provider, error) {
    // Validação da API key
    apiKey, ok := config["api_key"]
    if !ok {
        return nil, fmt.Errorf("api_key não encontrada na configuração do Claude")
    }
    
    // Configurações padrão
    model := "anthropic/claude-3-haiku"
    maxTokens := 4096
    temperature := 0.7
    baseURL := "https://openrouter.ai/api/v1"
    
    // Aplicar overrides de configuração
    if configModel, ok := config["model"]; ok && configModel != "" {
        model = configModel
    }
    // ... outros parâmetros
    
    return &ClaudeProvider{
        apiKey:      apiKey,
        model:       model,
        maxTokens:   maxTokens,
        temperature: temperature,
        baseURL:     baseURL,
    }, nil
}
```

### Métodos da Interface

```go
// Identificação do provider
func (p *ClaudeProvider) Name() string {
    return "Claude"
}

// Geração de conteúdo
func (p *ClaudeProvider) GenerateContent(prompt string) (string, error) {
    // Preparar requisição
    req := claudeRequest{
        Model: p.model,
        Messages: []claudeMessage{
            {Role: "user", Content: prompt},
        },
        MaxTokens:   p.maxTokens,
        Temperature: p.temperature,
    }
    
    // Enviar requisição HTTP
    // ... lógica de requisição
    
    return response.Choices[0].Message.Content, nil
}
```

## Configuração

### Variáveis de Ambiente

O Claude Provider suporta as seguintes variáveis de ambiente:

```bash
# Obrigatória
CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"

# Opcional - Provider padrão
ZION_AI_PROVIDER="claude"

# Opcional - Configurações do modelo
CLAUDE_MODEL="anthropic/claude-3-haiku"
CLAUDE_MAX_TOKENS="4096"
CLAUDE_TEMPERATURE="0.7"
CLAUDE_BASE_URL="https://openrouter.ai/api/v1"
```

### Configuração no config.yaml

```yaml
claude_api_key: ""
claude_model: "anthropic/claude-3-haiku"
claude_max_tokens: "4096"
claude_temperature: "0.7"
ai_provider: "claude"  # Para usar Claude como padrão
```

## Funções de Configuração

### config/config.go

```go
type Config struct {
    GeminiAPIKey     string
    OpenAIAPIKey     string
    OpenRouterAPIKey string
    ClaudeAPIKey     string  // ✅ Novo campo
    AIProvider       string
    HomeDir          string
    PluginsDir       string
}

func LoadConfig() *Config {
    // Detecção automática de provider
    aiProvider := os.Getenv("ZION_AI_PROVIDER")
    if aiProvider == "" {
        if os.Getenv("GEMINI_API_KEY") != "" {
            aiProvider = "gemini"
        } else if os.Getenv("OPENAI_API_KEY") != "" {
            aiProvider = "gpt"
        } else if os.Getenv("OPENROUTER_API_KEY") != "" {
            aiProvider = "openrouter"
        } else if os.Getenv("CLAUDE_API_KEY") != "" {
            aiProvider = "claude"  // ✅ Novo
        } else {
            aiProvider = "gemini"
        }
    }
    
    return &Config{
        GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
        OpenAIAPIKey:     os.Getenv("OPENAI_API_KEY"),
        OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
        ClaudeAPIKey:     os.Getenv("CLAUDE_API_KEY"),  // ✅ Novo
        AIProvider:       aiProvider,
        HomeDir:          zionDir,
        PluginsDir:       pluginsDir,
    }
}

func (c *Config) GetAIConfig() map[string]string {
    switch c.AIProvider {
    case "claude":  // ✅ Novo caso
        return map[string]string{
            "api_key":     c.ClaudeAPIKey,
            "model":       os.Getenv("CLAUDE_MODEL"),
            "max_tokens":  os.Getenv("CLAUDE_MAX_TOKENS"),
            "temperature": os.Getenv("CLAUDE_TEMPERATURE"),
            "base_url":    os.Getenv("CLAUDE_BASE_URL"),
        }
    // ... outros casos
    }
}
```

## Uso Prático

### Comandos CLI

```bash
# Listar providers (incluindo Claude)
zion provider list

# Configurar Claude
zion provider config claude

# Testar conexão Claude
zion provider test claude

# Gerar projeto usando Claude
zion scaffold -l python -n meu-projeto -d "Descrição do projeto" -p claude -k "sk-or-v1-sua-chave"

# Gerar projeto com modelo específico
zion scaffold -l typescript -n web-app -d "App web moderno" -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-sonnet"
```

### Exemplos de Uso

```bash
# Configuração via variáveis de ambiente
export CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"
export ZION_AI_PROVIDER="claude"
export CLAUDE_MODEL="anthropic/claude-3-haiku"

# Uso simples
zion scaffold -l javascript -n hello-world -d "Aplicação Hello World"

# Uso com parâmetros específicos
zion scaffold -l go -n api-server -d "API REST com autenticação" \
  -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-sonnet"
```

## Modelos Suportados

### Família Claude 3

- **claude-3-haiku**: Rápido e econômico
- **claude-3-sonnet**: Equilibrado (velocidade/qualidade)
- **claude-3-opus**: Máxima qualidade
- **claude-3-5-sonnet**: Última geração

### Formato de Modelo

```
anthropic/claude-3-haiku
anthropic/claude-3-sonnet
anthropic/claude-3-opus
anthropic/claude-3-5-sonnet
```

## Tratamento de Erros

### Erros Comuns

1. **API Key não configurada**
   ```
   api_key não encontrada na configuração do Claude
   ```

2. **Modelo não encontrado**
   ```
   API retornou status 400: Model not found
   ```

3. **Problema de autenticação**
   ```
   API retornou status 401: Unauthorized
   ```

4. **Limite de rate**
   ```
   API retornou status 429: Too Many Requests
   ```

### Soluções

```bash
# Verificar configuração
zion provider current

# Testar conexão
zion provider test claude

# Usar modelo padrão
zion scaffold -l python -n teste -p claude -k "sua-chave" -m "anthropic/claude-3-haiku"
```

## Integração com OpenRouter

O Claude Provider utiliza o OpenRouter como backend, que oferece:

- ✅ Acesso unificado aos modelos Claude
- ✅ Gestão de rate limiting
- ✅ Monitoramento de uso
- ✅ Suporte a múltiplos modelos

### Comparação: Claude vs OpenRouter

| Aspecto | Claude Provider | OpenRouter Provider |
|---------|-----------------|-------------------|
| **Uso** | Específico para Claude | Múltiplos modelos |
| **Configuração** | Simplificada | Mais flexível |
| **Modelos** | Apenas Claude | Claude + outros |
| **Interface** | Dedicada | Genérica |

## Roadmap

### Funcionalidades Futuras

- [ ] Suporte a Claude API direta (sem OpenRouter)
- [ ] Configuração avançada de parâmetros
- [ ] Cache de respostas
- [ ] Métricas de uso
- [ ] Suporte a streaming

### Melhorias Planejadas

- [ ] Otimização de tokens
- [ ] Tratamento de erros mais detalhado
- [ ] Suporte a context window maior
- [ ] Integração com ferramentas Claude

## Contribuição

Para contribuir com o Claude Provider:

1. Implemente melhorias no arquivo `ai/providers/claude.go`
2. Adicione testes unitários
3. Atualize a documentação
4. Envie um Pull Request

## Referências

- [OpenRouter API](https://openrouter.ai/docs)
- [Claude API Documentation](https://docs.anthropic.com/claude/reference)
- [Zion CLI GitHub](https://github.com/ktfth/zion)
