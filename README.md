# 🚀 Zion CLI

Uma ferramenta de scaffolding moderna que utiliza IA para gerar estruturas de projetos em qualquer linguagem de programação, com foco em boas práticas e modularidade.

## ✨ Características

- 🤖 **Integração com IA** - Utiliza Gemini AI para gerar estruturas inteligentes
- 🔌 **Sistema de Plugins** - Arquitetura extensível através de plugins
- 🌍 **Multi-linguagem** - Suporte a qualquer linguagem de programação
- 📦 **Scaffolding Inteligente** - Estruturas otimizadas para cada tipo de projeto
- 🛠️ **Configurável** - Personalizável através de plugins e configurações
- 🪟 **Suporte a Windows** - Funciona perfeitamente em ambientes Windows

## 🚦 Começando

### Pré-requisitos

- Go 1.20 ou superior
- Chave de API do Gemini (gratuita)

### Instalação

Você tem duas opções para instalar o Zion:

#### 1. Via go install (Recomendado)

```bash
go install github.com/ktfh/zion@latest

# Configure o ambiente (cria diretórios necessários e plugin de exemplo)
zion setup
```

#### 2. Compilando do Fonte

```bash
# Clone o repositório
git clone https://github.com/ktfth/zion.git

# Entre no diretório
cd zion

# Compile o projeto
go build

# Configure o ambiente (cria diretórios necessários e plugin de exemplo)
./zion setup
```

### Configuração

1. Obtenha uma chave de API do Gemini em: https://makersuite.google.com/app/apikey
2. Configure a chave de API:
   ```bash
   # Windows PowerShell
   $env:GEMINI_API_KEY="sua-chave-aqui"
   
   # Linux/macOS
   export GEMINI_API_KEY="sua-chave-aqui"
   ```

## 📚 Uso

### Gerar um novo projeto

```bash
# Sintaxe básica
zion scaffold -l <linguagem> -n <nome-projeto> -d "<descrição>"

# Exemplo: Criar uma API GraphQL em TypeScript
zion scaffold -l typescript -n minha-api -d "API GraphQL com autenticação e banco de dados"
```

### Comandos Disponíveis

- `zion setup` - Configura o ambiente inicial
- `zion scaffold` - Gera um novo projeto
  - `-l, --language` - Linguagem do projeto
  - `-n, --name` - Nome do projeto
  - `-d, --description` - Descrição do projeto

## 🔌 Sistema de Plugins

O Zion possui um sistema de plugins robusto que permite estender suas funcionalidades:

### Tipos de Plugins

1. **Plugins Estáticos** (Windows)
   - Implementados diretamente no código fonte
   - Requerem recompilação do Zion

2. **Plugins Dinâmicos** (Linux/macOS)
   - Arquivos `.so` carregados em tempo de execução
   - Podem ser adicionados sem recompilação

### Hooks Disponíveis

- `BeforeGeneration` - Executado antes da geração do scaffold
- `ModifyPrompt` - Permite modificar o prompt enviado à IA
- `AfterGeneration` - Executado após a geração do scaffold

### Diretório de Plugins

- Windows: `%APPDATA%\Local\Zion\plugins`
- Linux/macOS: `~/.zion/plugins`

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add: nova funcionalidade'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 🆘 Suporte

- **Issues**: Use o GitHub Issues para reportar problemas
- **Documentação**: Consulte a [Wiki](link-para-wiki) para documentação detalhada
- **Exemplos**: Veja a pasta `examples/` para exemplos de uso

## ⚠️ Notas Importantes

1. No Windows, os plugins são implementados estaticamente devido a limitações do Go com plugins dinâmicos no Windows
2. A chave API do Gemini é necessária para o funcionamento da ferramenta
3. Alguns caracteres especiais (como @ em pacotes npm) podem requerer tratamento especial

---
⭐️ Se este projeto te ajudou, considere dar uma estrela! 
