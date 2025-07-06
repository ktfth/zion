# Resumo das Implementações dos Providers

## Estrutura Consistente

Todos os providers seguem a mesma estrutura:

### 1. Interface Provider
```go
type Provider interface {
    Name() string
    GenerateContent(prompt string) (string, error)
}
```

### 2. Estrutura do Provider
- Struct contendo configurações (API Key, modelo, tokens, etc.)
- Função construtora `NewXXXProvider(config map[string]string) (Provider, error)`
- Método `Name() string`
- Método `GenerateContent(prompt string) (string, error)`

### 3. Configurações Suportadas

#### Gemini Provider
- `api_key` (obrigatório)
- Modelo fixo: `gemini-2.0-flash`
- Timeout: 60 segundos

#### GPT Provider
- `api_key` (obrigatório)
- `model` (opcional, padrão: `gpt-3.5-turbo`)
- `max_tokens` (opcional, padrão: `2048`)
- `temperature` (opcional, padrão: `0.7`)
- Timeout: 60 segundos

#### OpenRouter Provider
- `api_key` (obrigatório)
- `model` (opcional, padrão: `meta-llama/llama-3.2-3b-instruct:free`)
- `max_tokens` (opcional, padrão: `2048`)
- `temperature` (opcional, padrão: `0.7`)
- `base_url` (opcional, padrão: `https://openrouter.ai/api/v1`)
- Timeout: 60 segundos
- Headers customizados: HTTP-Referer e X-Title

## Consistências Implementadas

1. **Timeouts**: Todos os providers usam 60 segundos de timeout
2. **Validação de configuração**: Verificação de strings vazias (`!= ""`)
3. **Tratamento de erros**: Mensagens de erro detalhadas
4. **Estrutura de resposta**: Validação de estruturas vazias
5. **Headers**: Content-Type: application/json em todos

## Variáveis de Ambiente Suportadas

### Gemini
- `GEMINI_API_KEY`

### GPT/OpenAI
- `OPENAI_API_KEY`
- `OPENAI_MODEL`
- `OPENAI_MAX_TOKENS`
- `OPENAI_TEMPERATURE`

### OpenRouter
- `OPENROUTER_API_KEY`
- `OPENROUTER_MODEL`
- `OPENROUTER_MAX_TOKENS`
- `OPENROUTER_TEMPERATURE`
- `OPENROUTER_BASE_URL`

### Geral
- `ZION_AI_PROVIDER` (define qual provider usar)

## Registro dos Providers

Todos os providers são registrados automaticamente no `ProviderManager`:

```go
// No NewProviderManager()
manager.RegisterFactory("gemini", NewGeminiProvider)
manager.RegisterFactory("gpt", NewGPTProvider)
manager.RegisterFactory("openrouter", NewOpenRouterProvider)

// No init()
DefaultManager.RegisterFactory("gemini", NewGeminiProvider)
DefaultManager.RegisterFactory("gpt", NewGPTProvider)
DefaultManager.RegisterFactory("openrouter", NewOpenRouterProvider)
```

## Seleção Automática do Provider

O sistema seleciona automaticamente o provider baseado na disponibilidade das chaves:

1. Se `ZION_AI_PROVIDER` estiver definido, usa o especificado
2. Senão, verifica as chaves disponíveis na ordem:
   - `GEMINI_API_KEY` → usa "gemini"
   - `OPENAI_API_KEY` → usa "gpt"
   - `OPENROUTER_API_KEY` → usa "openrouter"
   - Padrão: "gemini"

## Uso

```go
// Obter configuração
cfg := config.LoadConfig()
aiConfig := cfg.GetAIConfig()

// Criar provider
provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, aiConfig)
if err != nil {
    log.Fatal(err)
}

// Usar provider
response, err := provider.GenerateContent("Seu prompt aqui")
if err != nil {
    log.Fatal(err)
}
```

Todas as implementações estão agora consistentes e seguem o mesmo padrão de design!
