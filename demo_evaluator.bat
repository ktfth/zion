@echo off
REM Script de demonstração do sistema Evaluator do Zion

echo 🔍 DEMONSTRAÇÃO DO SISTEMA EVALUATOR
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo.

REM Compilar o projeto
echo 🔨 Compilando o Zion...
go build -o zion.exe .
if %errorlevel% neq 0 (
    echo ❌ Erro na compilação!
    exit /b 1
)
echo ✅ Compilação concluída!
echo.

REM Teste 1: Projeto com boa qualidade
echo 📊 TESTE 1: Avaliando projeto com BOA qualidade
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
zion.exe evaluate -f example_project_good.json -l go --details
echo.

REM Teste 2: Projeto com problemas
echo 📊 TESTE 2: Avaliando projeto com PROBLEMAS
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
zion.exe evaluate -f example_project_bad.json -l go --details
echo.

REM Teste 3: JSON output
echo 📄 TESTE 3: Output em formato JSON
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
zion.exe evaluate -f example_project_good.json -l go --format json
echo.

echo ✨ Demonstração concluída!
echo.
echo 💡 Comandos disponíveis:
echo    • zion evaluate -f ^<arquivo^> -l ^<linguagem^>      # Avaliar projeto
echo    • zion scaffold -l ^<lang^> -n ^<nome^> -d ^<desc^>    # Criar com avaliação
echo    • zion scaffold --skip-evaluation ...            # Pular avaliação
echo.
pause
