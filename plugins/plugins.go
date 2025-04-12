package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"zion/config"
)

// ScaffoldHook define os pontos de extensão para plugins durante a geração de scaffold
type ScaffoldHook string

const (
	// BeforeGeneration é executado antes da geração do scaffold
	BeforeGeneration ScaffoldHook = "before_generation"
	// AfterGeneration é executado após a geração do scaffold
	AfterGeneration ScaffoldHook = "after_generation"
	// ModifyPrompt permite modificar o prompt antes de enviá-lo para a API
	ModifyPrompt ScaffoldHook = "modify_prompt"
)

// Plugin define a interface que todo plugin deve implementar.
type Plugin interface {
	// Name retorna o nome do plugin.
	Name() string
	// Execute contém a lógica a ser executada pelo plugin.
	Execute() error
	// GetHooks retorna os hooks que o plugin implementa
	GetHooks() map[ScaffoldHook]interface{}
}

// ScaffoldContext contém informações sobre o processo de geração de scaffold
type ScaffoldContext struct {
	ProjectName string
	Language    string
	Description string
	Prompt      string
	Response    string
}

// Mapa que mantém os plugins registrados.
var registeredPlugins = make(map[string]Plugin)

// RegisterPlugin permite o registro de um plugin.
func RegisterPlugin(p Plugin) {
	registeredPlugins[p.Name()] = p
	fmt.Printf("Plugin registrado: %s\n", p.Name())
}

// ListPlugins retorna os nomes dos plugins registrados.
func ListPlugins() []string {
	var names []string
	for name := range registeredPlugins {
		names = append(names, name)
	}
	return names
}

// ExecutePlugins executa a função Execute de cada plugin registrado.
func ExecutePlugins() {
	for name, plugin := range registeredPlugins {
		fmt.Println("Executando plugin:", name)
		if err := plugin.Execute(); err != nil {
			fmt.Printf("Erro na execução do plugin %s: %v\n", name, err)
		}
	}
}

// LoadPlugins carrega todos os plugins do diretório de plugins
func LoadPlugins() error {
	cfg := config.LoadConfig()
	if cfg.PluginsDir == "" {
		return fmt.Errorf("diretório de plugins não configurado")
	}

	// Normaliza o caminho do diretório de plugins
	pluginsDir := filepath.Clean(cfg.PluginsDir)

	// Verifica se o diretório existe
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		return fmt.Errorf("diretório de plugins não encontrado: %s", pluginsDir)
	}

	// Lista todos os arquivos no diretório
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de plugins: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !isPluginFile(entry.Name()) {
			continue
		}

		pluginPath := filepath.Join(pluginsDir, entry.Name())
		p, err := loadPlugin(pluginPath)
		if err != nil {
			fmt.Printf("⚠️  Aviso: erro ao carregar plugin %s: %v\n", entry.Name(), err)
			continue
		}

		registeredPlugins[p.Name()] = p
		fmt.Printf("✅ Plugin carregado: %s\n", p.Name())
	}

	return nil
}

// isPluginFile verifica se o arquivo é um plugin válido
func isPluginFile(filename string) bool {
	ext := filepath.Ext(filename)
	return strings.EqualFold(ext, ".so") || strings.EqualFold(ext, ".dll")
}

// loadPlugin carrega um único plugin
func loadPlugin(path string) (Plugin, error) {
	// Abre o arquivo do plugin
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir plugin: %v", err)
	}

	// Procura pelo símbolo "Plugin"
	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("símbolo 'Plugin' não encontrado: %v", err)
	}

	// Converte para a interface Plugin
	p, ok := symPlugin.(Plugin)
	if !ok {
		return nil, fmt.Errorf("símbolo 'Plugin' não implementa a interface Plugin")
	}

	return p, nil
}

// ExecuteHook executa todos os plugins que implementam um determinado hook
func ExecuteHook(hook ScaffoldHook, ctx *ScaffoldContext) *ScaffoldContext {
	for name, plugin := range registeredPlugins {
		hooks := plugin.GetHooks()
		if hookFunc, exists := hooks[hook]; exists {
			fmt.Printf("Executando hook %s do plugin %s\n", hook, name)

			switch hook {
			case BeforeGeneration:
				if beforeFunc, ok := hookFunc.(func(*ScaffoldContext) error); ok {
					if err := beforeFunc(ctx); err != nil {
						fmt.Printf("Erro na execução do hook %s do plugin %s: %v\n", hook, name, err)
					}
				}
			case AfterGeneration:
				if afterFunc, ok := hookFunc.(func(*ScaffoldContext) error); ok {
					if err := afterFunc(ctx); err != nil {
						fmt.Printf("Erro na execução do hook %s do plugin %s: %v\n", hook, name, err)
					}
				}
			case ModifyPrompt:
				if modifyFunc, ok := hookFunc.(func(string) string); ok {
					ctx.Prompt = modifyFunc(ctx.Prompt)
				}
			}
		}
	}

	return ctx
}
