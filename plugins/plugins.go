package plugins

import (
	"fmt"
	"os"
	"runtime"

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

// LoadPlugins carrega os plugins
func LoadPlugins(cfg *config.Config) error {
	// No Windows, os plugins são carregados estaticamente
	if runtime.GOOS == "windows" {
		fmt.Println("Plugins dinâmicos não são suportados no Windows. Usando plugins estáticos.")
		// Os plugins estáticos já são carregados automaticamente via init()
		return nil
	}

	// Em sistemas Unix, carregamos plugins dinâmicos do diretório de plugins
	fmt.Printf("Carregando plugins do diretório: %s\n", cfg.PluginsDir)

	// Verificar se o diretório de plugins existe
	if _, err := os.Stat(cfg.PluginsDir); os.IsNotExist(err) {
		fmt.Println("Diretório de plugins não encontrado, pulando carregamento de plugins")
		return nil
	}

	// Em sistemas Unix, implementaríamos o carregamento dinâmico de plugins
	// Mas como estamos focando na solução para Windows, isso fica como um TODO
	fmt.Println("Carregamento dinâmico de plugins ainda não implementado para sistemas Unix")

	return nil
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

