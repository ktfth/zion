package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateTypeScriptProject cria um projeto TypeScript básico com Express
// Esta função lida especificamente com o problema do @ nos nomes de pacotes npm
func CreateTypeScriptProject(projectName string) error {
	fmt.Println("Iniciando criação de projeto TypeScript:", projectName)
	
	// Criar o diretório raiz do projeto
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}

	fmt.Println("Criando estrutura do projeto em:", projectDir)

	// Criar estrutura de diretórios
	directories := []string{
		"src",
		"dist",
	}

	for _, dir := range directories {
		dirPath := filepath.Join(projectDir, dir)
		fmt.Println("Criando diretório:", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
	}

	// Criar arquivos
	files := map[string]string{
		"package.json": `{
  "name": "` + projectName + `",
  "version": "1.0.0",
  "description": "Simple Hello World API in TypeScript",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "nodemon src/index.ts"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "devDependencies": {
    "@types/express": "^4.17.17",
    "@types/node": "^20.4.5",
    "nodemon": "^3.0.1",
    "typescript": "^5.1.6"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}`,
		"tsconfig.json": `{
  "compilerOptions": {
    "target": "es6",
    "module": "commonjs",
    "rootDir": "./src",
    "outDir": "./dist",
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "skipLibCheck": true
  }
}`,
		"src/index.ts": `// Import the Express library
import express, { Request, Response } from 'express';

// Create a new Express application
const app = express();
const port = 3000;

// Define a route handler for the '/' endpoint
app.get('/', (req: Request, res: Response) => {
  res.send('Hello World!');
});

// Start the server and listen on the specified port
app.listen(port, () => {
  console.log("Server listening on port " + port);
});`,
		"README.md": `# Hello World API in TypeScript

A simple Hello World API built with TypeScript and Express.

## Prerequisites

- Node.js
- npm or yarn

## Installation

1. Clone the repository:
   ` + "```bash" + `
   git clone [repository URL]
   cd ` + projectName + `
   ` + "```" + `

2. Install dependencies:
   ` + "```bash" + `
   npm install
   # or
   yarn install
   ` + "```" + `

## Build

` + "```bash" + `
npm run build
# or
yarn build
` + "```" + `

## Run

` + "```bash" + `
npm start
# or
yarn start
` + "```" + `

The API will be accessible at http://localhost:3000.

## Development

Use the following command for development with automatic restart on file changes:

` + "```bash" + `
npm run dev
# or
yarn dev
` + "```" + ``,
		"hello.md": `# Welcome to the Hello World API!

This is a basic API that returns 'Hello World!' when you access the root route ('/').

Enjoy using it as a starting point for your TypeScript and Express projects!`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(projectDir, filePath)
		
		// Garantir que o diretório pai exista
		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório pai para %s: %v", filePath, err)
		}
		
		fmt.Println("Criando arquivo:", filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
		}
	}

	fmt.Println("\nEstrutura do projeto TypeScript criada com sucesso em:", projectDir)
	return nil
}
