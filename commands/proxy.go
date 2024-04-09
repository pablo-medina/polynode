package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"polynode/shared"
)

func SetProxyURL(httpProxy string) error {
	proxyConfigFile := filepath.Join(shared.GetInstallPath(), "proxy.json")

	// Crear la configuración del proxy
	proxyConfig := shared.ProxyConfig{
		HTTPProxy: httpProxy,
	}

	// Serializar la configuración del proxy a JSON
	proxyJSON, err := json.MarshalIndent(proxyConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("Error al serializar la configuración del proxy: %w", err)
	}

	// Escribir la configuración del proxy en el archivo proxy.json
	err = os.WriteFile(proxyConfigFile, proxyJSON, 0644)
	if err != nil {
		return fmt.Errorf("Error al escribir la configuración del proxy en el archivo: %w", err)
	}

	fmt.Printf("Se actualizó la configuración para usar el proxy '%s'\n", httpProxy)
	return nil
}
