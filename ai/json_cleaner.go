package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// CleanJSONString limpa e corrige problemas comuns em strings JSON retornadas pela API
func CleanJSONString(jsonStr string) string {
	// Primeiro, vamos verificar se o JSON já é válido
	var testObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &testObj); err == nil {
		// Se não houver erro, o JSON já é válido
		return jsonStr
	}

	// Remover caracteres de controle e espaços em branco extras
	cleaned := strings.TrimSpace(jsonStr)

	// Corrigir problema com pacotes npm
	cleaned = FixNpmPackagesInJSON(cleaned)

	// Verificar se a string limpa é um JSON válido
	if err := json.Unmarshal([]byte(cleaned), &testObj); err == nil {
		return cleaned
	}

	// Se ainda não for válido, tentar uma abordagem mais agressiva
	fmt.Println("Tentando abordagem mais agressiva para corrigir o JSON...")

	// Tentar remover caracteres problemáticos
	problematicChars := []string{"\u0000", "\u001F", "\u007F"}
	for _, char := range problematicChars {
		cleaned = strings.ReplaceAll(cleaned, char, "")
	}

	// Última tentativa: substituir aspas simples por aspas duplas em valores
	cleaned = FixQuotesInJSON(cleaned)

	return cleaned
}

// FixNpmPackagesInJSON corrige problemas com pacotes npm em JSON
func FixNpmPackagesInJSON(jsonStr string) string {
	// Problema comum: aspas em nomes de pacotes npm que começam com @
	// Exemplo: "@types/express": "^4.17.17"
	
	// Primeiro, vamos identificar o bloco de package.json
	packageJsonRegex := regexp.MustCompile(`"package\.json":\s*"\{([^}]+)\}"`) 
	matches := packageJsonRegex.FindStringSubmatch(jsonStr)
	
	if len(matches) > 1 {
		// Encontramos o bloco de package.json
		packageJsonContent := matches[1]
		
		// Corrigir aspas escapadas dentro do package.json
		fixed := packageJsonContent
		
		// Substituir \" por "
		fixed = strings.ReplaceAll(fixed, "\\\"", "\"")
		
		// Substituir o conteúdo original pelo corrigido
		jsonStr = strings.Replace(jsonStr, packageJsonContent, fixed, 1)
	}
	
	// Procurar por padrões específicos que causam problemas
	// Problema com o @types/express e similares
	// Busca por padrões como: "@types/express": "^4.17.17",
	problemRegex := regexp.MustCompile(`"(@[^"]+)":\s*"([^"]+)",\s*\n\s*"([^"]+)"`) 
	problemMatches := problemRegex.FindAllStringSubmatch(jsonStr, -1)
	
	for _, match := range problemMatches {
		if len(match) > 3 {
			original := match[0]
			packageName := match[1]  // @types/express
			version := match[2]      // ^4.17.17
			nextPackage := match[3]  // typescript ou outro pacote
			
			// Criar a versão corrigida
			corrected := fmt.Sprintf(`"%s": "%s",
    "%s"`, packageName, version, nextPackage)
			
			// Substituir na string JSON
			jsonStr = strings.Replace(jsonStr, original, corrected, 1)
		}
	}
	
	return jsonStr
}

// FixQuotesInJSON corrige problemas com aspas em JSON
func FixQuotesInJSON(jsonStr string) string {
	// Problema comum: aspas simples em vez de aspas duplas em valores
	// Exemplo: "dependencies": { 'express': '^4.18.2' }
	
	// Substituir aspas simples por aspas duplas apenas em valores
	// Isso é uma solução simplificada e pode não funcionar em todos os casos
	valueRegex := regexp.MustCompile(`:\s*'([^']*)'`)
	fixed := valueRegex.ReplaceAllString(jsonStr, `: "$1"`)
	
	return fixed
}
