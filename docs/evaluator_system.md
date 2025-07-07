# Sistema de Evaluator para Zion

O sistema de evaluator foi implementado para verificar a qualidade e aderência a boas práticas dos projetos antes de materializá-los.

## Funcionalidades

### 1. **Avaliação Automática**
- Integrada ao comando `scaffold` por padrão
- Verifica qualidade antes de materializar o projeto
- Bloqueia criação de projetos com issues críticos

### 2. **Avaliação Manual**
- Comando `evaluate` para análise independente
- Suporte a múltiplos formatos de entrada
- Relatórios detalhados com scores e sugestões

### 3. **Regras de Avaliação**
- **Estrutura de Diretórios**: Verifica convenções por linguagem
- **Nomenclatura de Arquivos**: Valida padrões de naming
- **Arquivos Obrigatórios**: Verifica presença de arquivos essenciais
- **Dependências**: Analisa consistência e vulnerabilidades
- **Configurações**: Valida arquivos de build e configuração
- **Melhores Práticas**: Verifica documentação, testes, etc.
- **Segurança**: Identifica vulnerabilidades conhecidas

## Uso

### Avaliação Automática (Integrada ao Scaffold)
```bash
# Avaliação automática habilitada por padrão
zion scaffold -l go -n meu-projeto -d "API REST com PostgreSQL"

# Pular avaliação se necessário
zion scaffold -l go -n meu-projeto -d "API REST" --skip-evaluation
```

### Avaliação Manual
```bash
# Avaliar arquivo de resposta da IA
zion evaluate -f response.txt -l go

# Formato JSON para automação
zion evaluate -f project.json -l typescript --format json

# Relatório detalhado
zion evaluate -f structure.json -l python --details
```

## Critérios de Qualidade

### Scores
- **90-100**: Excelente 🏆
- **75-89**: Bom ✅
- **60-74**: Regular ⚡
- **40-59**: Ruim ⚠️
- **0-39**: Crítico ❌

### Severidades de Issues
- **Critical**: Bloqueia materialização do projeto
- **High**: Problemas importantes que devem ser corrigidos
- **Medium**: Melhorias recomendadas
- **Low**: Sugestões de otimização
- **Info**: Informações gerais

## Categorias de Análise

### 1. Estrutura (Structure)
- Organização de diretórios
- Convenções por linguagem
- Arquivos na localização correta

### 2. Nomenclatura (Naming)
- Convenções de nomenclatura
- Consistência de nomes
- Padrões por linguagem

### 3. Dependências (Dependencies)
- Versões apropriadas
- Dependências duplicadas
- Vulnerabilidades conhecidas

### 4. Configuração (Configuration)
- Arquivos de build válidos
- Configurações recomendadas
- Compatibilidade de versões

### 5. Segurança (Security)
- Vulnerabilidades em dependências
- Arquivos sensíveis expostos
- Configurações inseguras

### 6. Manutenibilidade (Maintainability)
- Documentação adequada
- Estrutura de testes
- Melhores práticas

## Exemplos de Regras por Linguagem

### Go
- ✅ Estrutura padrão: `cmd/`, `pkg/`, `internal/`
- ✅ Arquivo `go.mod` presente
- ✅ Nomenclatura em `snake_case`
- ✅ Makefile para automação

### JavaScript/TypeScript
- ✅ Diretório `src/` ou `lib/`
- ✅ Arquivo `package.json` válido
- ✅ Scripts de build configurados
- ✅ TypeScript em modo strict

### Python
- ✅ Estrutura de pacote com `__init__.py`
- ✅ `requirements.txt` ou `pyproject.toml`
- ✅ Nomenclatura em `snake_case`
- ✅ Configuração de testes

### Java
- ✅ Estrutura Maven/Gradle
- ✅ `pom.xml` ou `build.gradle`
- ✅ Nomenclatura PascalCase para classes
- ✅ Wrappers de build

## Personalização

### Adicionando Novas Regras
```go
// Implementar interface EvaluationRule
type MinhaRegra struct{}

func (r *MinhaRegra) Name() string { return "MinhaRegra" }
func (r *MinhaRegra) Description() string { return "Descrição da regra" }
func (r *MinhaRegra) Category() Category { return CategoryCustom }
func (r *MinhaRegra) Weight() float64 { return 5.0 }
func (r *MinhaRegra) Evaluate(structure *ProjectStructure) (float64, []Issue) {
    // Lógica de avaliação
    return score, issues
}

// Registrar a regra
evaluator.RegisterRule(&MinhaRegra{})
```

### Configuração de Pesos
- Estrutura: 15 pontos
- Arquivos Obrigatórios: 12 pontos
- Dependências: 10 pontos
- Melhores Práticas: 10 pontos
- Nomenclatura: 8 pontos
- Configuração: 8 pontos
- Testes: 8 pontos
- Build: 7 pontos
- Documentação: 6 pontos

## Benefícios

### Para Desenvolvedores
- ✅ Projetos seguem melhores práticas desde o início
- ✅ Feedback imediato sobre qualidade
- ✅ Aprendizado de convenções por linguagem
- ✅ Detecção precoce de problemas de segurança

### Para Equipes
- ✅ Consistência entre projetos
- ✅ Padrões de qualidade automáticos
- ✅ Redução de code review para questões básicas
- ✅ Documentação automática de práticas

### Para Projetos
- ✅ Maior manutenibilidade
- ✅ Menos vulnerabilidades
- ✅ Melhor documentação
- ✅ Estrutura organizacional clara

## Próximos Passos

1. **Extensão de Regras**: Adicionar mais regras específicas por framework
2. **Integração com IDEs**: Plugin para VS Code com feedback em tempo real
3. **Cache de Avaliações**: Evitar re-avaliações desnecessárias
4. **Relatórios HTML**: Interface web para relatórios detalhados
5. **Métricas de Tendência**: Acompanhar evolução da qualidade ao longo do tempo

O sistema de evaluator garante que os projetos gerados pelo Zion sigam as melhores práticas desde o primeiro commit, elevando a qualidade geral do desenvolvimento.
