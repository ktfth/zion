# Funcionalidade llms.txt - Scaffolding Contextualizado

## Visão Geral

A funcionalidade `llms.txt` permite que o Zion CLI leia informações contextuais sobre o projeto atual para gerar scaffolding mais preciso e personalizado. Esta funcionalidade é especialmente útil para:

- Expansão de projetos existentes
- Adição de recursos a códigos já escritos
- Manutenção da consistência com a arquitetura atual
- Geração contextualizada baseada em especificações

## Como Funciona

### 1. Detecção Automática
Quando o comando `zion scaffold` é executado, a ferramenta automaticamente:
- Verifica se existe um arquivo `llms.txt` no diretório atual
- Se encontrado, entra automaticamente em "modo contextual"
- Lê e processa o conteúdo do arquivo
- Analisa a estrutura existente do projeto

### 2. Processamento do llms.txt
O arquivo `llms.txt` é processado para extrair:
- **Seções temáticas** (marcadas com # ou ## ou ###)
- **Referências de arquivos** (caminhos iniciados com ./ ou /)
- **Descrição do projeto**
- **Requisitos técnicos**
- **Arquitetura e tecnologias**

### 3. Enriquecimento do Prompt
O prompt enviado para a IA é enriquecido com:
- Contexto completo do projeto
- Estrutura de arquivos existente
- Conteúdo dos arquivos referenciados
- Especificações técnicas
- Padrões e convenções já estabelecidos

## Formato do Arquivo llms.txt

### Estrutura Básica
```
# Título do Projeto

## Descrição
Descrição geral do projeto e seus objetivos.

## Tecnologias
Lista das tecnologias utilizadas.

## Arquitetura
Explicação da arquitetura do projeto.

## Requisitos
Requisitos específicos para novas funcionalidades.

## Arquivos de Contexto
./src/main.js
./package.json
./README.md
```

### Seções Reconhecidas
- `objetivo`, `purpose` - Objetivo do projeto
- `descrição`, `description` - Descrição geral
- `tecnologias`, `tech stack` - Stack tecnológico
- `requisitos`, `requirements` - Requisitos técnicos
- `arquitetura`, `architecture` - Arquitetura do sistema

### Referências de Arquivos
- Caminhos relativos: `./src/components/Button.js`
- Caminhos absolutos: `/home/user/project/config.yaml`
- Os arquivos referenciados são lidos e incluídos no contexto

## Modos de Operação

### 1. Modo Automático
```bash
# Se llms.txt existe, modo contextual é ativado automaticamente
zion scaffold -d "Adicionar sistema de autenticação"
```

### 2. Modo Forçado
```bash
# Força modo contextual mesmo sem todos os parâmetros
zion scaffold --contextual -d "Adicionar API REST"
```

### 3. Detecção Automática de Linguagem
Se a linguagem não for especificada, o sistema tenta detectar baseado em:
- Arquivos característicos (package.json, go.mod, requirements.txt, etc.)
- Conteúdo do llms.txt
- Estrutura de diretórios

## Comportamento Inteligente

### 1. Merge de Arquivos
Quando arquivos já existem, o sistema decide como proceder:

- **Markdown (.md)**: Merge automático (adiciona ao final)
- **package.json**: Merge inteligente (combina dependências)
- **Configuração**: Backup e substituição
- **README**: Backup e substituição
- **Código fonte**: Skip (não sobrescreve)

### 2. Backup Automático
Arquivos importantes são automaticamente salvos com backup:
```
arquivo.json → arquivo.json.backup.1638360000
```

### 3. Detecção de Conflitos
O sistema detecta e evita:
- Sobrescrever arquivos importantes sem backup
- Conflitos de estrutura de diretórios
- Inconsistências de dependências

## Exemplos de Uso

### 1. Expandir Projeto JavaScript
```bash
# Dentro de um projeto com package.json
echo "# Meu App React
## Requisitos
Adicionar sistema de autenticação JWT
## Tecnologias
React, Node.js, Express, JWT
## Arquivos de Contexto
./package.json
./src/App.js" > llms.txt

zion scaffold -d "Implementar autenticação JWT"
```

### 2. Adicionar Features a Projeto Go
```bash
# Dentro de um projeto Go
echo "# API Go
## Objetivo
API REST com PostgreSQL
## Requisitos
Adicionar endpoints de usuário
## Arquivos de Contexto
./go.mod
./main.go" > llms.txt

zion scaffold -d "Endpoints CRUD para usuários"
```

### 3. Projeto Python com Contexto Complexo
```bash
echo "# ML Project
## Descrição
Projeto de Machine Learning para classificação
## Tecnologias
Python, scikit-learn, pandas, FastAPI
## Requisitos
Adicionar pipeline de preprocessing
## Arquitetura
- data/ - Dados de treinamento
- models/ - Modelos treinados  
- api/ - API FastAPI
- preprocessing/ - Scripts de preprocessing
## Arquivos de Contexto
./requirements.txt
./src/train.py
./config/settings.py" > llms.txt

zion scaffold -d "Pipeline de preprocessing com validação"
```

## Vantagens

### 1. Contextualização Precisa
- A IA entende o projeto existente
- Mantém consistência com padrões estabelecidos
- Evita conflitos de arquitetura

### 2. Workflow Fluido
- Não precisa recriar projeto do zero
- Expande funcionalidades incrementalmente
- Preserva trabalho já realizado

### 3. Inteligência Adaptativa
- Detecta automaticamente tecnologias
- Adapta-se ao estilo de código existente
- Sugere melhorias baseadas no contexto

### 4. Segurança
- Backup automático de arquivos importantes
- Não sobrescreve código sem aviso
- Permite revisão antes de aplicar mudanças

## Limitações e Considerações

### 1. Tamanho do Contexto
- Arquivos muito grandes são truncados (limite: 2000 caracteres)
- Muitos arquivos podem tornar o prompt muito grande
- Recomenda-se referenciar apenas arquivos essenciais

### 2. Tipos de Arquivo Suportados
- Merge inteligente limitado a tipos conhecidos
- Arquivos binários são ignorados
- Alguns formatos podem necessitar merge manual

### 3. Dependências
- Requer configuração de provider de IA
- Qualidade depende do modelo utilizado
- Pode necessitar múltiplas iterações para resultados ideais

## Troubleshooting

### Problema: Linguagem não detectada
**Solução**: Especifique manualmente com `-l linguagem`

### Problema: Arquivos não sendo mesclados corretamente
**Solução**: Use `--skip-evaluation` e revise manualmente

### Problema: Contexto muito grande
**Solução**: Reduza número de arquivos referenciados no llms.txt

### Problema: Backup não criado
**Solução**: Verifique permissões de escrita no diretório

## Próximos Passos

Esta funcionalidade abre caminho para:
- Templates contextuais personalizados
- Integração com IDEs
- Análise semântica de código
- Sugestões automáticas de melhorias
- Workflow de desenvolvimento assistido por IA
