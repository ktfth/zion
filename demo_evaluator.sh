#!/bin/bash

# Script de demonstração do sistema Evaluator do Zion

echo "🔍 DEMONSTRAÇÃO DO SISTEMA EVALUATOR"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Compilar o projeto
echo "🔨 Compilando o Zion..."
go build -o zion.exe .
if [ $? -ne 0 ]; then
    echo "❌ Erro na compilação!"
    exit 1
fi
echo "✅ Compilação concluída!"
echo ""

# Teste 1: Projeto com boa qualidade
echo "📊 TESTE 1: Avaliando projeto com BOA qualidade"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./zion.exe evaluate -f example_project_good.json -l go --details
echo ""

# Teste 2: Projeto com problemas
echo "📊 TESTE 2: Avaliando projeto com PROBLEMAS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./zion.exe evaluate -f example_project_bad.json -l go --details
echo ""

# Teste 3: Scaffold com avaliação automática (projeto bom)
echo "🚀 TESTE 3: Scaffold com avaliação automática (projeto bom)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Simulando scaffold de projeto bem estruturado..."
echo "Este teste criaria um projeto real com a estrutura avaliada positivamente."
echo ""

# Teste 4: Scaffold com avaliação automática (projeto problemático)
echo "⚠️  TESTE 4: Scaffold com avaliação automática (projeto problemático)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Simulando scaffold de projeto com problemas..."
echo "Este teste seria bloqueado pela avaliação automática."
echo ""

# Teste 5: JSON output
echo "📄 TESTE 5: Output em formato JSON"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./zion.exe evaluate -f example_project_good.json -l go --format json
echo ""

echo "✨ Demonstração concluída!"
echo ""
echo "💡 Comandos disponíveis:"
echo "   • zion evaluate -f <arquivo> -l <linguagem>      # Avaliar projeto"
echo "   • zion scaffold -l <lang> -n <nome> -d <desc>    # Criar com avaliação"
echo "   • zion scaffold --skip-evaluation ...            # Pular avaliação"
echo ""
