# Sistema de Geração em Camadas - Zion AI

## Visão Geral

O sistema de geração em camadas foi implementado para resolver problemas de overflow de contexto que podem ocorrer ao gerar projetos grandes com APIs de IA. Quando o contexto excede os limites do modelo (por exemplo, 200.000 tokens no OpenRouter), o sistema automaticamente divide a geração em múltiplas etapas.

## Como Funciona

### 1. Detecção de Overflow

O sistema monitora automaticamente o tamanho do contexto antes de enviar para a API:

```go
// Verifica se o prompt causaria overflow
if DetectContextOverflow(prompt, provider.Name()) {
    // Usar geração em camadas
    layeredGen, err := NewLayeredGenerator(language, projectName, description, llmsContext)
    // ...
}
```

### 2. Planejamento de Camadas

Quando detecta overflow, o sistema primeiro planeja as camadas necessárias:

- **Core**: Estrutura básica e configuração
- **Business**: Lógica de negócio e modelos  
- **API**: Endpoints e controladores
- **Frontend**: Componentes de interface (se aplicável)
- **Tests**: Testes unitários e integração
- **Docs**: Documentação

### 3. Geração Sequencial

Cada camada é gerada sequencialmente, mantendo contexto das camadas anteriores:

```go
for i, layerPlan := range layers {
    layerResult, err := lg.generateLayer(layerPlan, response.Layers)
    response.Layers = append(response.Layers, *layerResult)
}
```

### 4. Materialização

O projeto final é materializado combinando todas as camadas:

```go
CreateLayeredProject(projectName, layeredResponse)
```

## Limites de Tokens por Provider

| Provider   | Limite Seguro | Limite Real |
|------------|---------------|-------------|
| OpenRouter | 150.000       | 200.000     |
| GPT-4      | 100.000       | 128.000     |
| Claude     | 150.000       | 200.000     |
| Gemini 2.0 | 180.000       | 2.000.000   |

## Detecção Automática de Erros

O sistema detecta automaticamente erros de contexto através de padrões:

- "context length"
- "maximum context" 
- "token limit"
- "too many tokens"
- "input too long"

## Exemplo de Uso

```bash
# O sistema detecta automaticamente quando é necessário usar camadas
zion scaffold -l go -n meu-projeto-grande -d "Sistema complexo com microserviços, API GraphQL, autenticação JWT, banco de dados PostgreSQL, Redis cache, monitoramento com Prometheus, logs estruturados, CI/CD completo, testes unitários e integração"

# Saída esperada:
# ⚠️  Contexto muito grande detectado - usando geração em camadas
# 📋 Planejadas 4 camadas de geração
# 🔧 Gerando camada 1/4: core...
# ✅ Camada core concluída (5 arquivos, 3 diretórios)
# 🔧 Gerando camada 2/4: business...
# ...
```

## Estrutura de Saída

### Resposta em Camadas

```json
{
  "project_info": {
    "name": "meu-projeto",
    "language": "go", 
    "description": "Descrição do projeto"
  },
  "layers": [
    {
      "layer_name": "core",
      "description": "Estrutura básica do projeto",
      "directories": ["cmd", "pkg", "internal"],
      "files": {
        "main.go": {
          "content": "package main..."
        },
        "go.mod": {
          "content": "module meu-projeto..."
        }
      },
      "dependencies": ["go 1.21"],
      "next_steps": ["go mod tidy", "go run main.go"]
    }
  ]
}
```

### Arquivo de Resumo

O sistema gera automaticamente um arquivo `ZION_LAYERS_SUMMARY.md` com:

- Informações do projeto
- Lista de todas as camadas criadas
- Arquivos e diretórios por camada
- Dependências identificadas  
- Próximos passos sugeridos

## Vantagens

1. **Escalabilidade**: Suporta projetos de qualquer tamanho
2. **Confiabilidade**: Reduz falhas por overflow de contexto
3. **Organização**: Estrutura o projeto em camadas lógicas
4. **Transparência**: Fornece visibilidade total do processo
5. **Flexibilidade**: Adapta-se a diferentes linguagens e tipos de projeto

## Fallback Automático

Se a geração em camadas falhar, o sistema automaticamente:

1. Tenta usar o método tradicional de geração
2. Salva a resposta bruta para análise manual
3. Fornece instruções para resolução manual

## Monitoramento

Durante a execução, o sistema fornece feedback detalhado:

```
🏗️ Iniciando geração em camadas (limite: 150000 tokens)
📋 Planejadas 4 camadas de geração
🔧 Gerando camada 1/4: core...
✅ Camada core concluída (5 arquivos, 3 diretórios)
📁 Criando 8 diretórios...
⚙️  Materializando camada 1/4: core
🎉 Projeto criado com sucesso!
```

Este sistema garante que mesmo projetos muito complexos possam ser gerados consistentemente, respeitando os limites das APIs de IA e mantendo alta qualidade na saída.
