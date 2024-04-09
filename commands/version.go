package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"polynode/shared"
	"strings"
)

func ShowCurrentNodeVersion() {
	// Verificar si el directorio current existe
	_, err := os.Stat(shared.GetCurrentVersionPath())
	if err != nil && os.IsNotExist(err) {
		fmt.Println("No hay ninguna versión seleccionada. Utilice el comando use para seleccionar una.")
		return
	}

	// Obtener la ruta del archivo version.info
	versionFile := filepath.Join(shared.GetCurrentVersionPath(), "version.info")

	// Verificar si el archivo version.info existe
	_, err = os.Stat(versionFile)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("No se encontró version.info en la versión actual. Posiblemente fue instalada manualmente.")
		return
	}

	// Leer el contenido del archivo version.info
	versionBytes, err := os.ReadFile(versionFile)
	if err != nil {
		fmt.Println("Error al leer la versión actual:", err)
	}

	fmt.Println(strings.TrimSpace(string(versionBytes)))
}
