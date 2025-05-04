package ai

// DirectParseAndCreateProject é uma função que faz o parsing direto da resposta da API
// e cria a estrutura do projeto sem depender do parsing JSON padrão
func DirectParseAndCreateProject(projectName, response string) error {
	// Primeiro, tentar o método padrão
	err := CreateProjectFromRawJSON(projectName, response)
	if err == nil {
		return nil
	}

	// Se falhar, tentar o método alternativo
	return ExtractAndCreateProject(projectName, response)
}

// ExtractDirectoriesAndFiles extrai diretórios e arquivos diretamente da resposta da API
// sem depender do parsing JSON padrão
func ExtractDirectoriesAndFiles(response string) ([]string, map[string]string) {
	directories := []string{}
	files := make(map[string]string)
	return directories, files
}
