param(
    [Parameter(Mandatory=$true)]
    [string]$ProjectName
)

Write-Host "Iniciando criação de projeto TypeScript: $ProjectName"

# Criar o diretório raiz do projeto
$projectDir = Join-Path -Path "." -ChildPath $ProjectName
New-Item -Path $projectDir -ItemType Directory -Force | Out-Null

Write-Host "Criando estrutura do projeto em: $projectDir"

# Criar estrutura de diretórios
$directories = @(
    "src",
    "dist"
)

foreach ($dir in $directories) {
    $dirPath = Join-Path -Path $projectDir -ChildPath $dir
    Write-Host "Criando diretório: $dir"
    New-Item -Path $dirPath -ItemType Directory -Force | Out-Null
}

# Criar arquivos
$packageJson = @"
{
  "name": "$ProjectName",
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
}
"@

$tsconfigJson = @"
{
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
}
"@

$indexTs = @"
// Import the Express library
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
});
"@

$readmeMd = @"
# Hello World API in TypeScript

A simple Hello World API built with TypeScript and Express.

## Prerequisites

- Node.js
- npm or yarn

## Installation

1. Clone the repository:
   ```bash
   git clone [repository URL]
   cd $ProjectName
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   yarn install
   ```

## Build

```bash
npm run build
# or
yarn build
```

## Run

```bash
npm start
# or
yarn start
```

The API will be accessible at http://localhost:3000.

## Development

Use the following command for development with automatic restart on file changes:

```bash
npm run dev
# or
yarn dev
```
"@

$helloMd = @"
# Welcome to the Hello World API!

This is a basic API that returns 'Hello World!' when you access the root route ('/').

Enjoy using it as a starting point for your TypeScript and Express projects!
"@

# Criar arquivos
$files = @{
    "package.json" = $packageJson
    "tsconfig.json" = $tsconfigJson
    "src/index.ts" = $indexTs
    "README.md" = $readmeMd
    "hello.md" = $helloMd
}

foreach ($file in $files.Keys) {
    $filePath = Join-Path -Path $projectDir -ChildPath $file
    $parentDir = Split-Path -Path $filePath -Parent
    
    # Garantir que o diretório pai exista
    if (-not (Test-Path -Path $parentDir)) {
        New-Item -Path $parentDir -ItemType Directory -Force | Out-Null
    }
    
    Write-Host "Criando arquivo: $file"
    Set-Content -Path $filePath -Value $files[$file]
}

Write-Host "`nEstrutura do projeto TypeScript criada com sucesso em: $projectDir"
