package ai

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ProcessNpmPackageJson processa especificamente respostas JSON que contêm package.json com @types
// Esta função é uma solução específica para o problema de parsing do JSON quando há pacotes npm com @
func ProcessNpmPackageJson(jsonString string) ScaffoldResponse {
	scaffoldResp := ScaffoldResponse{}
	scaffoldResp.Structure.Files = make(map[string]interface{}) // Alterar esta linha

	// Lógica para processar o JSON e popular scaffoldResp.Structure.Files
	// Ao invés de:
	// scaffoldResp.Structure.Files[key] = value.(string)
	// Use:
	// scaffoldResp.Structure.Files[key] = value
	// ou se precisar converter para string:
	// scaffoldResp.Structure.Files[key] = string(bytes)

	// Exemplo de como você pode estar processando o JSON:
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &data); err == nil {
		if files, ok := data["structure"].(map[string]interface{}); ok {
			scaffoldResp.Structure = struct {
				Directories []string                 `json:"directories"`
				Files       map[string]interface{} `json:"files"`
			}{
				Directories: []string{}, // Se necessário, extrair diretórios também
				Files:       files,
			}
		}
	}

	return scaffoldResp
}

// FixPackageJsonContent corrige o conteúdo do package.json para lidar com o problema do @
func FixPackageJsonContent(content string) string {
	// Primeiro, vamos verificar se o conteúdo já é um JSON válido
	var testObj interface{}
	if err := json.Unmarshal([]byte(content), &testObj); err == nil {
		// Se não houver erro, o JSON já é válido
		return content
	}
	
	// Corrige o problema específico com @types
	// Busca por linhas como: "@types/express": "^4.17.17",
	atTypesRegex := regexp.MustCompile(`"(@[^"]+)":\s*"([^"]+)",`)
	fixed := atTypesRegex.ReplaceAllString(content, "\"$1\": \"$2\",")
	
	// Corrige outros problemas comuns
	fixed = strings.ReplaceAll(fixed, "\n,", ",")
	fixed = strings.ReplaceAll(fixed, ",\n,", ",")
	
	return fixed
}
