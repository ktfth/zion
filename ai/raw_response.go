package ai

import (
	"fmt"
	"os"
	"path/filepath"
)

// SaveRawResponse salva a resposta bruta da API em um arquivo JSON
// para que o usuário possa processá-la manualmente posteriormente
func SaveRawResponse(projectName, response string) error {
	// Criar o diretório do projeto se não existir
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}

	// Caminho para o arquivo de resposta
	responsePath := filepath.Join(projectName, "zion_response.json")
	
	// Salvar a resposta bruta no arquivo
	if err := os.WriteFile(responsePath, []byte(response), 0644); err != nil {
		return fmt.Errorf("erro ao salvar resposta bruta: %v", err)
	}

	// Criar um README explicando como usar a resposta
	readmePath := filepath.Join(projectName, "README.md")
	readmeContent := `# Projeto gerado pelo Zion

Este projeto foi gerado pelo Zion, mas devido a limitações no parsing do JSON com caracteres especiais como '@',
a resposta bruta da API foi salva no arquivo 'zion_response.json'.

## Como proceder

1. Examine o arquivo 'zion_response.json' para ver a estrutura do projeto gerada pela IA.
2. Crie manualmente os arquivos e diretórios conforme especificado na resposta.
3. Para projetos TypeScript com pacotes npm que contêm '@', você pode iniciar um projeto com:

` + "```bash" + `
npm init -y
npm install express
npm install --save-dev typescript @types/express @types/node ts-node-dev
` + "```" + `

4. Crie um arquivo tsconfig.json básico:

` + "```json" + `
{
  "compilerOptions": {
    "target": "es2016",
    "module": "commonjs",
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "skipLibCheck": true,
    "outDir": "dist"
  }
}
` + "```" + `

5. Crie a estrutura de diretórios e arquivos conforme especificado na resposta.

## Próximos passos

Em versões futuras do Zion, este processo será automatizado para lidar corretamente com caracteres especiais.
`

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("erro ao criar README: %v", err)
	}

	fmt.Println("Resposta bruta salva em:", responsePath)
	fmt.Println("README com instruções criado em:", readmePath)
	
	return nil
}
