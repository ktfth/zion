# Como usar o OpenRouter com Zion CLI

## Configuração via Variáveis de Ambiente

### 1. Configurar as variáveis de ambiente:

```bash
# Definir a API key do OpenRouter
export OPENROUTER_API_KEY="sk-or-v1-your-api-key-here"

# Definir o OpenRouter como provider padrão
export ZION_AI_PROVIDER="openrouter"

# (Opcional) Definir um modelo específico
export OPENROUTER_MODEL="anthropic/claude-3-haiku"

# (Opcional) Configurar outros parâmetros
export OPENROUTER_MAX_TOKENS="4096"
export OPENROUTER_TEMPERATURE="0.8"
```

### 2. Verificar configuração:

```bash
# Listar providers disponíveis
zion provider list

# Ver configuração atual
zion provider current

# Testar conexão com OpenRouter
zion provider test openrouter
```

### 3. Usar o OpenRouter para gerar projetos:

```bash
# Usar configuração padrão
zion scaffold -l javascript -n meu-projeto -d "Um projeto React moderno"

# Usar configuração específica
zion scaffold -l python -n api-rest -d "API REST com FastAPI"
```

## Configuração via Linha de Comando

### 1. Usar OpenRouter diretamente na linha de comando:

```bash
# Usar OpenRouter com API key específica
zion scaffold \
  -l go \
  -n meu-app \
  -d "API REST em Go com Gin" \
  -p openrouter \
  -k "sk-or-v1-your-api-key-here"

# Usar OpenRouter com modelo específico
zion scaffold \
  -l typescript \
  -n frontend-app \
  -d "Aplicação React com TypeScript" \
  -p openrouter \
  -k "sk-or-v1-your-api-key-here" \
  -m "anthropic/claude-3-haiku"
```

### 2. Usar modelo gratuito do OpenRouter:

```bash
zion scaffold \
  -l python \
  -n ml-project \
  -d "Projeto de Machine Learning" \
  -p openrouter \
  -k "sk-or-v1-your-api-key-here" \
  -m "meta-llama/llama-3.2-3b-instruct:free"
```

## Comandos Úteis

### Verificar configuração:
```bash
# Listar todos os providers
zion provider list

# Ver configuração do OpenRouter
zion provider config openrouter

# Testar conexão
zion provider test openrouter

# Ver provider atual
zion provider current
```

### Obter ajuda:
```bash
# Ajuda geral
zion --help

# Ajuda do comando provider
zion provider --help

# Ajuda do comando scaffold
zion scaffold --help
```

## Modelos Recomendados

### Gratuitos:
- `meta-llama/llama-3.2-3b-instruct:free`
- `qwen/qwen-2-7b-instruct:free`

### Pagos (alta qualidade):
- `anthropic/claude-3-haiku`
- `openai/gpt-3.5-turbo`
- `google/gemini-pro`
- `meta-llama/llama-3.1-8b-instruct`

## Exemplo Completo

```bash
# 1. Configurar ambiente
export OPENROUTER_API_KEY="sk-or-v1-your-api-key-here"
export ZION_AI_PROVIDER="openrouter"
export OPENROUTER_MODEL="anthropic/claude-3-haiku"

# 2. Testar configuração
zion provider test openrouter

# 3. Gerar projeto
zion scaffold \
  -l typescript \
  -n ecommerce-app \
  -d "Aplicação de e-commerce com React, TypeScript e Node.js"

# 4. Verificar resultado
cd ecommerce-app
ls -la
```

O OpenRouter está agora totalmente integrado na CLI do Zion! 🚀
