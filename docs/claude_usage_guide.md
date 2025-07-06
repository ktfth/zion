# 🧠 Guia de Uso - Claude Provider

## Introdução

O Claude Provider oferece uma interface dedicada para usar os modelos Claude (Anthropic) no Zion CLI. Este guia cobre todas as funcionalidades e melhores práticas.

## Instalação e Configuração

### Pré-requisitos

1. **Conta no OpenRouter**: https://openrouter.ai/
2. **API Key do OpenRouter**: https://openrouter.ai/keys
3. **Créditos**: Para modelos Claude pagos

### Configuração Básica

```bash
# 1. Definir API Key
export CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"

# 2. Definir Claude como provider padrão
export ZION_AI_PROVIDER="claude"

# 3. Configurar modelo (opcional)
export CLAUDE_MODEL="anthropic/claude-3-haiku"
```

### Configuração Avançada

```bash
# Configurações completas
export CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"
export ZION_AI_PROVIDER="claude"
export CLAUDE_MODEL="anthropic/claude-3-sonnet"
export CLAUDE_MAX_TOKENS="4096"
export CLAUDE_TEMPERATURE="0.7"
export CLAUDE_BASE_URL="https://openrouter.ai/api/v1"
```

## Comandos Disponíveis

### Gerenciamento de Provider

```bash
# Listar todos os providers
zion provider list

# Ver provider atual
zion provider current

# Instruções de configuração
zion provider config claude

# Testar conexão
zion provider test claude
```

### Geração de Projetos

```bash
# Sintaxe básica
zion scaffold -l <linguagem> -n <nome> -d "<descrição>"

# Com provider específico
zion scaffold -l python -n meu-projeto -d "Descrição" -p claude

# Com API key específica
zion scaffold -l javascript -n app -d "Descrição" -p claude -k "sk-or-v1-sua-chave"

# Com modelo específico
zion scaffold -l typescript -n web-app -d "Descrição" -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-sonnet"
```

## Modelos Disponíveis

### Claude 3 Haiku (Recomendado para Início)

```bash
# Rápido e econômico
export CLAUDE_MODEL="anthropic/claude-3-haiku"

# Exemplo de uso
zion scaffold -l python -n api-rest -d "API REST com FastAPI" \
  -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-haiku"
```

**Características:**
- ⚡ Muito rápido
- 💰 Econômico
- 🎯 Ideal para projetos simples
- 📝 Boa documentação

### Claude 3 Sonnet (Equilibrado)

```bash
# Equilibrio entre velocidade e qualidade
export CLAUDE_MODEL="anthropic/claude-3-sonnet"

# Exemplo de uso
zion scaffold -l typescript -n fullstack-app -d "Aplicação completa com React e Node.js" \
  -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-sonnet"
```

**Características:**
- ⚖️ Equilibrado
- 🎨 Boa criatividade
- 🔧 Ideal para projetos médios
- 📊 Análise detalhada

### Claude 3 Opus (Máxima Qualidade)

```bash
# Máxima qualidade
export CLAUDE_MODEL="anthropic/claude-3-opus"

# Exemplo de uso
zion scaffold -l go -n microservice -d "Microserviço complexo com gRPC, autenticação e monitoramento" \
  -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-opus"
```

**Características:**
- 🏆 Máxima qualidade
- 🧠 Raciocínio complexo
- 💎 Ideal para projetos enterprise
- 📈 Melhor análise de requisitos

### Claude 3.5 Sonnet (Última Geração)

```bash
# Última geração
export CLAUDE_MODEL="anthropic/claude-3-5-sonnet"

# Exemplo de uso
zion scaffold -l rust -n blockchain-app -d "Aplicação blockchain com smart contracts" \
  -p claude -k "sk-or-v1-sua-chave" -m "anthropic/claude-3-5-sonnet"
```

**Características:**
- 🆕 Última geração
- 🚀 Melhor performance
- 💡 Insights avançados
- 🔬 Ideal para projetos inovadores

## Exemplos Práticos

### Exemplo 1: Aplicação Web Básica

```bash
# Configurar ambiente
export CLAUDE_API_KEY="sk-or-v1-sua-chave-openrouter"
export ZION_AI_PROVIDER="claude"
export CLAUDE_MODEL="anthropic/claude-3-haiku"

# Gerar projeto
zion scaffold -l javascript -n blog-app -d "Blog pessoal com React e Node.js"

# Resultado esperado
cd blog-app
ls -la
# README.md, package.json, src/, public/, etc.
```

### Exemplo 2: API REST Enterprise

```bash
# Configurar para projeto complexo
export CLAUDE_MODEL="anthropic/claude-3-sonnet"

# Gerar API avançada
zion scaffold -l python -n enterprise-api -d "API REST enterprise com autenticação JWT, cache Redis, banco PostgreSQL, monitoramento com Prometheus, logs estruturados e documentação automática"

# Verificar estrutura
cd enterprise-api
tree
```

### Exemplo 3: Aplicação Mobile

```bash
# Gerar app mobile
zion scaffold -l typescript -n mobile-app -d "Aplicação mobile multiplataforma com React Native, navegação stack, gerenciamento de estado com Redux, integração com APIs REST e push notifications"

# Configurar dependências
cd mobile-app
npm install
```

### Exemplo 4: Projeto DevOps

```bash
# Usar Claude para infraestrutura
export CLAUDE_MODEL="anthropic/claude-3-opus"

zion scaffold -l yaml -n k8s-manifests -d "Manifests Kubernetes para aplicação web com deployment, services, ingress, configmaps, secrets, HPA e monitoramento"

# Verificar manifests
cd k8s-manifests
kubectl apply --dry-run=client -f .
```

## Comparação com OpenRouter

### Usando Claude via OpenRouter

```bash
# OpenRouter com modelo Claude
export ZION_AI_PROVIDER="openrouter"
export OPENROUTER_API_KEY="sk-or-v1-sua-chave"
export OPENROUTER_MODEL="anthropic/claude-3-haiku"

zion scaffold -l python -n test-openrouter -d "Teste via OpenRouter"
```

### Usando Claude Provider Dedicado

```bash
# Provider Claude dedicado
export ZION_AI_PROVIDER="claude"
export CLAUDE_API_KEY="sk-or-v1-sua-chave"
export CLAUDE_MODEL="anthropic/claude-3-haiku"

zion scaffold -l python -n test-claude -d "Teste via Claude Provider"
```

### Diferenças Práticas

| Aspecto | OpenRouter | Claude Provider |
|---------|------------|-----------------|
| **Comando** | `zion scaffold -p openrouter -m "anthropic/claude-3-haiku"` | `zion scaffold -p claude -m "anthropic/claude-3-haiku"` |
| **Configuração** | `OPENROUTER_MODEL` | `CLAUDE_MODEL` |
| **Flexibilidade** | Acesso a múltiplos modelos | Focado em Claude |
| **Simplicidade** | Mais genérico | Mais específico |

## Troubleshooting

### Problema: Provider não encontrado

```bash
# Erro
❌ Erro: provedor não suportado: claude

# Solução
go build -o zion.exe .
zion provider list
```

### Problema: API Key inválida

```bash
# Erro
❌ Erro na geração: API retornou status 401

# Verificar configuração
echo $CLAUDE_API_KEY
zion provider test claude

# Testar com OpenRouter
export OPENROUTER_API_KEY="$CLAUDE_API_KEY"
export ZION_AI_PROVIDER="openrouter"
export OPENROUTER_MODEL="anthropic/claude-3-haiku"
zion provider test openrouter
```

### Problema: Modelo não encontrado

```bash
# Erro
❌ Erro: Model not found

# Usar modelo válido
zion scaffold -l python -n teste -p claude -k "sua-chave" -m "anthropic/claude-3-haiku"
```

### Problema: Rate Limiting

```bash
# Erro
❌ Erro: API retornou status 429

# Aguardar e tentar novamente
sleep 60
zion scaffold -l python -n teste -d "Teste após rate limit"
```

## Otimização de Uso

### Escolha do Modelo por Projeto

```bash
# Prototipagem rápida
export CLAUDE_MODEL="anthropic/claude-3-haiku"

# Projetos comerciais
export CLAUDE_MODEL="anthropic/claude-3-sonnet"

# Projetos enterprise
export CLAUDE_MODEL="anthropic/claude-3-opus"

# Projetos experimentais
export CLAUDE_MODEL="anthropic/claude-3-5-sonnet"
```

### Configuração de Tokens

```bash
# Projetos simples
export CLAUDE_MAX_TOKENS="2048"

# Projetos médios
export CLAUDE_MAX_TOKENS="4096"

# Projetos complexos
export CLAUDE_MAX_TOKENS="8192"
```

### Ajuste de Temperatura

```bash
# Código estruturado (menos criativo)
export CLAUDE_TEMPERATURE="0.3"

# Padrão (equilibrado)
export CLAUDE_TEMPERATURE="0.7"

# Código criativo (mais variação)
export CLAUDE_TEMPERATURE="0.9"
```

## Scripts de Automação

### Script de Configuração

```bash
#!/bin/bash
# setup-claude.sh

# Verificar se API key foi fornecida
if [ -z "$1" ]; then
    echo "Uso: $0 <api-key>"
    echo "Exemplo: $0 sk-or-v1-sua-chave-aqui"
    exit 1
fi

# Configurar Claude
export CLAUDE_API_KEY="$1"
export ZION_AI_PROVIDER="claude"
export CLAUDE_MODEL="anthropic/claude-3-haiku"
export CLAUDE_MAX_TOKENS="4096"
export CLAUDE_TEMPERATURE="0.7"

# Testar configuração
echo "Testando configuração Claude..."
zion provider test claude

if [ $? -eq 0 ]; then
    echo "✅ Claude configurado com sucesso!"
    echo "Modelo: $CLAUDE_MODEL"
    echo "Tokens: $CLAUDE_MAX_TOKENS"
    echo "Temperature: $CLAUDE_TEMPERATURE"
else
    echo "❌ Erro na configuração do Claude"
    exit 1
fi
```

### Script de Geração de Projetos

```bash
#!/bin/bash
# generate-project.sh

LANGUAGE="$1"
PROJECT_NAME="$2"
DESCRIPTION="$3"
MODEL="${4:-anthropic/claude-3-haiku}"

if [ -z "$LANGUAGE" ] || [ -z "$PROJECT_NAME" ] || [ -z "$DESCRIPTION" ]; then
    echo "Uso: $0 <linguagem> <nome-projeto> <descrição> [modelo]"
    echo "Exemplo: $0 python minha-api 'API REST com FastAPI' anthropic/claude-3-sonnet"
    exit 1
fi

echo "Gerando projeto: $PROJECT_NAME"
echo "Linguagem: $LANGUAGE"
echo "Descrição: $DESCRIPTION"
echo "Modelo: $MODEL"
echo "---"

zion scaffold -l "$LANGUAGE" -n "$PROJECT_NAME" -d "$DESCRIPTION" -p claude -m "$MODEL"

if [ $? -eq 0 ]; then
    echo "✅ Projeto gerado com sucesso!"
    echo "📁 Diretório: $PROJECT_NAME"
    echo "💡 Para começar: cd $PROJECT_NAME"
else
    echo "❌ Erro na geração do projeto"
    exit 1
fi
```

## Integração com IDEs

### VS Code

```json
// .vscode/settings.json
{
    "terminal.integrated.env.linux": {
        "CLAUDE_API_KEY": "sk-or-v1-sua-chave-aqui",
        "ZION_AI_PROVIDER": "claude",
        "CLAUDE_MODEL": "anthropic/claude-3-haiku"
    }
}
```

### Configuração de Tarefas

```json
// .vscode/tasks.json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Zion: Gerar Projeto Claude",
            "type": "shell",
            "command": "zion",
            "args": [
                "scaffold",
                "-l", "${input:language}",
                "-n", "${input:projectName}",
                "-d", "${input:description}",
                "-p", "claude",
                "-m", "anthropic/claude-3-haiku"
            ],
            "group": "build",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            }
        }
    ],
    "inputs": [
        {
            "id": "language",
            "description": "Linguagem do projeto",
            "default": "javascript",
            "type": "promptString"
        },
        {
            "id": "projectName",
            "description": "Nome do projeto",
            "type": "promptString"
        },
        {
            "id": "description",
            "description": "Descrição do projeto",
            "type": "promptString"
        }
    ]
}
```

## Conclusão

O Claude Provider oferece uma interface especializada para usar os modelos Claude no Zion CLI, proporcionando:

- ✅ Configuração simplificada
- ✅ Interface dedicada
- ✅ Melhor experiência de uso
- ✅ Suporte a todos os modelos Claude
- ✅ Integração perfeita com o ecossistema Zion

Para dúvidas ou suporte, consulte a documentação oficial ou abra uma issue no GitHub.

---

**Próximos Passos:**
1. Configure sua API key do OpenRouter
2. Teste com o modelo Claude 3 Haiku
3. Explore diferentes tipos de projetos
4. Ajuste configurações conforme necessário
5. Contribua com melhorias!
