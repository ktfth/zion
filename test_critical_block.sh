#!/bin/bash
# Script de teste para verificar se o scaffold bloqueia projetos críticos

echo "🧪 Testando sistema de avaliação crítica do Zion"
echo "================================================"

# Criar um arquivo temporário que simula uma resposta da IA com problemas críticos
cat > temp_ai_response.txt << 'EOF'
{
  "structure": {
    "files": {
      ".env": {
        "type": "file",
        "content": "# Critical security issues\nAPI_KEY=super_secret_key_123\nDB_PASSWORD=admin123\nJWT_SECRET=my_secret_token\nPRIVATE_KEY=-----BEGIN PRIVATE KEY-----"
      },
      "config.js": {
        "type": "file",
        "content": "// Missing package.json for JS project\nmodule.exports = {\n  database: {\n    password: 'hardcoded_password',\n    secret: 'api_secret_123'\n  }\n};"
      },
      "src/index.js": {
        "type": "file",
        "content": "const password = 'admin123';\nconst apiKey = 'secret_key';\nconsole.log('Hello World');"
      }
    }
  }
}
EOF

echo ""
echo "📋 Avaliando projeto crítico:"
./zion evaluate -f temp_ai_response.txt -l javascript

echo ""
echo "🚨 Resultado esperado: Projeto deve ser bloqueado por issues críticos"
echo ""

# Limpar
rm -f temp_ai_response.txt

echo "✅ Teste concluído"
