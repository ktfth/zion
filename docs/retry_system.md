# Sistema de Retry do Zion CLI

## Visão Geral

O Zion CLI implementa um sistema inteligente de retry que melhora significativamente a confiabilidade da geração de projetos. Como os modelos de IA podem ocasionalmente falhar ou retornar respostas inconsistentes, o sistema de retry garante que o projeto seja criado com sucesso mesmo quando há falhas temporárias.

## Funcionalidades

### 1. Retry Automático
- **Padrão**: 3 tentativas automáticas
- **Configurável**: Pode ser ajustado via flag `--retries` ou `-r`
- **Inteligente**: Diferentes estratégias para diferentes tipos de falha

### 2. Dois Níveis de Retry

#### Nível 1: Geração de Conteúdo
- Retry na comunicação com a API de IA
- Intervalo de 2 segundos entre tentativas
- Mantém o contexto e configurações originais

#### Nível 2: Criação da Estrutura
- Retry na criação física dos arquivos e diretórios
- Intervalo de 1 segundo entre tentativas
- Fallback para salvamento da resposta bruta

### 3. Feedback Visual
- Indicadores de progresso para cada tentativa
- Contadores de tentativas na saída final
- Mensagens claras sobre o status de cada retry

## Uso

### Comando Básico
```bash
zion scaffold -l go -n meu-projeto -d "API REST com PostgreSQL"
```

### Configurando Número de Retries
```bash
# Definir 5 tentativas máximas
zion scaffold -l typescript -n webapp -d "App React" --retries 5

# Usando forma abreviada
zion scaffold -l python -n ml-project -d "Machine Learning" -r 2
```

### Exemplos Práticos

#### Projeto Go com Retry Padrão
```bash
zion scaffold -l go -n api-service -d "Microserviço com gRPC e PostgreSQL"
```

#### Projeto TypeScript com Retry Customizado
```bash
zion scaffold -l typescript -n frontend-app -d "SPA com React e Redux" -r 5
```

#### Projeto Python com Retry Mínimo
```bash
zion scaffold -l python -n data-pipeline -d "Pipeline de dados com Apache Airflow" -r 1
```

## Comportamento Detalhado

### Fluxo de Retry para Geração
1. **Tentativa 1**: Execução normal
2. **Tentativa 2-N**: Em caso de falha
   - Aguarda 2 segundos
   - Mantém todas as configurações originais
   - Reutiliza o mesmo prompt e contexto
   - Informa o número da tentativa

### Fluxo de Retry para Criação
1. **Tentativa 1**: Criação normal da estrutura
2. **Tentativa 2-N**: Em caso de falha
   - Aguarda 1 segundo
   - Tenta recriar a estrutura completa
   - Informa o número da tentativa

### Fallback Final
Se todas as tentativas de criação falharem:
- Salva a resposta bruta em `README.md`
- Permite análise manual da resposta
- Não interrompe o processo completamente

## Saída de Exemplo

### Sucesso na Primeira Tentativa
```
🤖 Gerando estrutura com IA... ✅
📂 Criando estrutura do projeto... ✅

✨ Projeto criado com sucesso! ✨
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📁 Local: meu-projeto
⏱️  Tempo total: 3.45 segundos
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Sucesso com Retry
```
🤖 Gerando estrutura com IA...
⚠️  Erro na tentativa 1: API rate limit exceeded
🔄 Tentando novamente em 2 segundos...
🔄 Tentativa 2/3... ✅
📂 Criando estrutura do projeto... ✅

✨ Projeto criado com sucesso! ✨
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📁 Local: meu-projeto
⏱️  Tempo total: 8.12 segundos
🔄 Tentativas de geração: 2/3
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Falha Total
```
🤖 Gerando estrutura com IA...
⚠️  Erro na tentativa 1: Connection timeout
🔄 Tentando novamente em 2 segundos...
🔄 Tentativa 2/3...
⚠️  Erro na tentativa 2: Connection timeout
🔄 Tentando novamente em 2 segundos...
🔄 Tentativa 3/3...
❌ Falha após 3 tentativas:
Connection timeout
```

## Casos de Uso

### Redes Instáveis
- Increase retries para conexões intermitentes
- Exemplo: `--retries 5` para ambientes com conectividade limitada

### APIs com Rate Limiting
- Retries padrão geralmente suficientes
- O intervalo de 2 segundos ajuda a evitar rate limits

### Desenvolvimento Local
- Pode usar `--retries 1` para feedback rápido
- Útil durante testes e desenvolvimento

### Produção/CI
- Manter padrão de 3 tentativas
- Balanceio entre confiabilidade e tempo de execução

## Benefícios

1. **Maior Confiabilidade**: Reduz falhas por problemas temporários
2. **Experiência do Usuário**: Processo mais fluido e previsível
3. **Flexibilidade**: Configurável conforme necessidade
4. **Transparência**: Feedback claro sobre o que está acontecendo
5. **Robustez**: Múltiplas estratégias de recuperação

## Considerações de Performance

- **Tempo Adicional**: Cada retry adiciona tempo de execução
- **Recursos**: Múltiplas tentativas consomem mais recursos da API
- **Balanceamento**: Padrão de 3 tentativas oferece bom equilíbrio

## Troubleshooting

### Muitas Falhas
- Verifique conexão com internet
- Confirme se as credenciais de API estão corretas
- Considere usar provider alternativo

### Tempo de Execução Longo
- Reduza número de retries com `-r 1`
- Verifique se a descrição do projeto não é muito complexa

### Estrutura Incompleta
- Verifique o arquivo `README.md` gerado
- Analise os logs de erro detalhados
- Considere simplificar a descrição do projeto

## Versão e Compatibilidade

- **Versão**: Implementado na versão atual do Zion CLI
- **Compatibilidade**: Funciona com todos os providers de IA suportados
- **Dependências**: Não requer bibliotecas adicionais
