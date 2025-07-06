# OpenRouter Provider Configuration

O provider OpenRouter permite usar vários modelos de IA através da API do OpenRouter.

## Configuração

### Variáveis de Ambiente

```bash
# Chave da API OpenRouter (obrigatória)
export OPENROUTER_API_KEY="sk-or-v1-your-api-key-here"

# Modelo a ser usado (opcional, padrão: meta-llama/llama-3.2-3b-instruct:free)
export OPENROUTER_MODEL="meta-llama/llama-3.2-3b-instruct:free"

# Número máximo de tokens (opcional, padrão: 2048)
export OPENROUTER_MAX_TOKENS="2048"

# Temperatura para controlar a criatividade (opcional, padrão: 0.7)
export OPENROUTER_TEMPERATURE="0.7"

# URL base customizada (opcional, padrão: https://openrouter.ai/api/v1)
export OPENROUTER_BASE_URL="https://openrouter.ai/api/v1"

# Definir o provider como OpenRouter
export ZION_AI_PROVIDER="openrouter"
```

### Arquivo config.yaml

```yaml
gemini_api_key: ""
openai_api_key: ""
openrouter_api_key: "sk-or-v1-your-api-key-here"
ai_provider: "openrouter"
openai_model: "gpt-3.5-turbo"
openai_max_tokens: "2048"
openai_temperature: "0.7"
openrouter_model: "meta-llama/llama-3.2-3b-instruct:free"
openrouter_max_tokens: "2048"
openrouter_temperature: "0.7"
home_dir: "C:\\Users\\Usuario\\.zion"
plugins_dir: "C:\\Users\\Usuario\\.zion\\plugins"
```

## Modelos Disponíveis

O OpenRouter oferece acesso a vários modelos, incluindo:

- `meta-llama/llama-3.2-3b-instruct:free` (gratuito)
- `openai/gpt-3.5-turbo`
- `openai/gpt-4`
- `anthropic/claude-3-haiku`
- `google/gemini-pro`
- E muitos outros...

## Características

- **Timeout**: 60 segundos para chamadas à API
- **Headers customizados**: Inclui referer e título para identificação
- **Configuração flexível**: Permite customizar modelo, tokens, temperatura e URL base
- **Compatível com OpenAI API**: Usa o mesmo formato de mensagens
- **Tratamento de erros**: Mensagens de erro detalhadas

## Exemplo de Uso

```go
import "github.com/zion-ai/providers"

// Configuração
config := map[string]string{
    "api_key": "sk-or-v1-your-api-key-here",
    "model": "meta-llama/llama-3.2-3b-instruct:free",
    "max_tokens": "2048",
    "temperature": "0.7",
}

// Criar provider
provider, err := NewOpenRouterProvider(config)
if err != nil {
    log.Fatal(err)
}

// Gerar conteúdo
response, err := provider.GenerateContent("Olá, como você está?")
if err != nil {
    log.Fatal(err)
}

fmt.Println(response)
```

## Obter Chave da API

1. Acesse [OpenRouter.ai](https://openrouter.ai)
2. Crie uma conta ou faça login
3. Vá para a seção "Keys"
4. Crie uma nova chave API
5. Use a chave no formato `sk-or-v1-...`
