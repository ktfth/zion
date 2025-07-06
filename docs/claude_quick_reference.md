# 🧠 Claude Provider - Referência Rápida

## Configuração Rápida

```bash
# Básico
export CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"
export ZION_AI_PROVIDER="claude"

# Completo
export CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"
export ZION_AI_PROVIDER="claude"
export CLAUDE_MODEL="anthropic/claude-3-haiku"
export CLAUDE_MAX_TOKENS="4096"
export CLAUDE_TEMPERATURE="0.7"
```

## Comandos Essenciais

```bash
# Gerenciar providers
zion provider list                    # Listar providers
zion provider current                 # Ver provider atual
zion provider config claude          # Configurar Claude
zion provider test claude            # Testar Claude

# Gerar projetos
zion scaffold -l js -n app -d "desc"                                    # Básico
zion scaffold -l py -n api -d "desc" -p claude                         # Provider específico
zion scaffold -l ts -n web -d "desc" -p claude -k "sua-chave"         # Com API key
zion scaffold -l go -n srv -d "desc" -p claude -m "claude-3-sonnet"   # Com modelo
```

## Modelos Disponíveis

| Modelo | Código | Uso | Velocidade | Qualidade | Custo |
|--------|--------|-----|------------|-----------|-------|
| Claude 3 Haiku | `anthropic/claude-3-haiku` | Geral | ⚡⚡⚡ | ⭐⭐⭐ | 💰 |
| Claude 3 Sonnet | `anthropic/claude-3-sonnet` | Médio | ⚡⚡ | ⭐⭐⭐⭐ | 💰💰 |
| Claude 3 Opus | `anthropic/claude-3-opus` | Enterprise | ⚡ | ⭐⭐⭐⭐⭐ | 💰💰💰 |
| Claude 3.5 Sonnet | `anthropic/claude-3-5-sonnet` | Última geração | ⚡⚡ | ⭐⭐⭐⭐⭐ | 💰💰💰 |

## Linguagens e Tipos de Projeto

### JavaScript/TypeScript
```bash
# React App
zion scaffold -l javascript -n react-app -d "Aplicação React moderna" -p claude

# Node.js API
zion scaffold -l typescript -n api-server -d "API REST com Express" -p claude

# Full Stack
zion scaffold -l typescript -n fullstack -d "App completa Next.js" -p claude
```

### Python
```bash
# FastAPI
zion scaffold -l python -n fastapi-app -d "API REST com FastAPI" -p claude

# Django
zion scaffold -l python -n django-app -d "Aplicação web Django" -p claude

# Data Science
zion scaffold -l python -n ml-project -d "Projeto ML com pandas/sklearn" -p claude
```

### Go
```bash
# API REST
zion scaffold -l go -n go-api -d "API REST em Go" -p claude

# gRPC Service
zion scaffold -l go -n grpc-service -d "Serviço gRPC" -p claude

# CLI Tool
zion scaffold -l go -n cli-tool -d "Ferramenta CLI" -p claude
```

### Outros
```bash
# Rust
zion scaffold -l rust -n rust-app -d "Aplicação Rust" -p claude

# C#
zion scaffold -l csharp -n dotnet-api -d "API .NET" -p claude

# Java
zion scaffold -l java -n spring-app -d "App Spring Boot" -p claude
```

## Variáveis de Ambiente

| Variável | Obrigatória | Padrão | Descrição |
|----------|-------------|---------|-----------|
| `CLAUDE_API_KEY` | ✅ | - | Chave API do OpenRouter |
| `ZION_AI_PROVIDER` | ❌ | auto | Provider padrão |
| `CLAUDE_MODEL` | ❌ | `claude-3-haiku` | Modelo Claude |
| `CLAUDE_MAX_TOKENS` | ❌ | `4096` | Tokens máximos |
| `CLAUDE_TEMPERATURE` | ❌ | `0.7` | Temperatura (0-1) |
| `CLAUDE_BASE_URL` | ❌ | OpenRouter URL | URL base da API |

## Troubleshooting

### Erros Comuns
```bash
# Provider não encontrado
❌ provedor não suportado: claude
✅ go build -o zion.exe . && zion provider list

# API key inválida
❌ API retornou status 401
✅ zion provider test claude

# Modelo não encontrado
❌ Model not found
✅ zion scaffold -p claude -m "anthropic/claude-3-haiku"

# Rate limiting
❌ API retornou status 429
✅ sleep 60 && zion scaffold ...
```

### Verificação de Status
```bash
# Verificar configuração
echo $CLAUDE_API_KEY
echo $ZION_AI_PROVIDER
echo $CLAUDE_MODEL

# Testar conexão
zion provider test claude

# Ver provider atual
zion provider current
```

## Comparação: Claude vs OpenRouter

### Claude Provider
```bash
export ZION_AI_PROVIDER="claude"
export CLAUDE_API_KEY="sk-or-v1-sua-chave"
export CLAUDE_MODEL="anthropic/claude-3-haiku"
zion scaffold -l python -n app -d "Descrição"
```

### OpenRouter com Claude
```bash
export ZION_AI_PROVIDER="openrouter"
export OPENROUTER_API_KEY="sk-or-v1-sua-chave"
export OPENROUTER_MODEL="anthropic/claude-3-haiku"
zion scaffold -l python -n app -d "Descrição"
```

| Aspecto | Claude Provider | OpenRouter |
|---------|-----------------|------------|
| **Foco** | Apenas Claude | Múltiplos modelos |
| **Configuração** | Específica | Genérica |
| **Simplicidade** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Flexibilidade** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

## Recomendações por Uso

### 🚀 Prototipagem (Rápida)
```bash
export CLAUDE_MODEL="anthropic/claude-3-haiku"
export CLAUDE_MAX_TOKENS="2048"
export CLAUDE_TEMPERATURE="0.7"
```

### 🏢 Projetos Comerciais
```bash
export CLAUDE_MODEL="anthropic/claude-3-sonnet"
export CLAUDE_MAX_TOKENS="4096"
export CLAUDE_TEMPERATURE="0.5"
```

### 🏛️ Enterprise
```bash
export CLAUDE_MODEL="anthropic/claude-3-opus"
export CLAUDE_MAX_TOKENS="8192"
export CLAUDE_TEMPERATURE="0.3"
```

### 🔬 Experimental
```bash
export CLAUDE_MODEL="anthropic/claude-3-5-sonnet"
export CLAUDE_MAX_TOKENS="8192"
export CLAUDE_TEMPERATURE="0.8"
```

## Scripts Úteis

### Configuração Automática
```bash
#!/bin/bash
# setup-claude-quick.sh
export CLAUDE_API_KEY="$1"
export ZION_AI_PROVIDER="claude"
export CLAUDE_MODEL="anthropic/claude-3-haiku"
zion provider test claude && echo "✅ Claude configurado!"
```

### Geração Rápida
```bash
#!/bin/bash
# quick-generate.sh
zion scaffold -l "$1" -n "$2" -d "$3" -p claude -m "anthropic/claude-3-haiku"
```

## Links Úteis

- [OpenRouter.ai](https://openrouter.ai/) - Plataforma
- [OpenRouter Keys](https://openrouter.ai/keys) - Chaves API
- [Claude Models](https://openrouter.ai/models?q=claude) - Modelos disponíveis
- [Pricing](https://openrouter.ai/pricing) - Preços
- [Zion Docs](../README.md) - Documentação principal

---

**Comandos de Emergência:**
```bash
# Reset total
unset CLAUDE_API_KEY ZION_AI_PROVIDER CLAUDE_MODEL
export ZION_AI_PROVIDER="gemini"  # Fallback gratuito

# Diagnóstico completo
zion provider list
zion provider current
zion provider test claude
```
