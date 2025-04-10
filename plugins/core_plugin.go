package plugins

import "fmt"

// CorePlugin é o plugin principal do sistema.
type CorePlugin struct{}

func (p CorePlugin) Name() string {
	return "CorePlugin"
}

func (p CorePlugin) Execute() error {
	fmt.Println("CorePlugin executado com sucesso!")
	return nil
}

// GetHooks retorna os hooks que o plugin implementa
func (p CorePlugin) GetHooks() map[ScaffoldHook]interface{} {
	// O CorePlugin não implementa nenhum hook específico
	return map[ScaffoldHook]interface{}{}
}

// A função init é chamada automaticamente e registra o plugin CorePlugin.
func init() {
	RegisterPlugin(CorePlugin{})
}
