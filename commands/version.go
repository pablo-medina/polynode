package commands

import (
	"fmt"
	"polynode/shared"
	"strings"
)

func ShowCurrentNodeVersion() {
	// Obtener versión actual
	currentVersion := shared.GetCurrentVersion()

	if currentVersion == "" {
		fmt.Println("No hay ninguna versión seleccionada. Utilice el comando use para seleccionar una.")
		return
	}

	fmt.Println(strings.TrimSpace(string(currentVersion)))
}
