package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"polynode/shared"
)

func UninstallNodeVersion(version string) error {
	versionDir := filepath.Join(shared.GetRepoPath(), "node-v"+version+"-win-x64")

	// Verificar si la versión que se intenta desinstalar está instalada
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("La versión %s no está instalada", version)
	}

	// Eliminar el directorio de la versión
	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("Error al desinstalar la versión %s: %w", version, err)
	}

	// Si la versión desinstalada es la misma que la actual, borrar el directorio "current"
	currentVersion := shared.GetCurrentVersion()
	if currentVersion == version {
		if err := os.RemoveAll(shared.GetCurrentVersionPath()); err != nil {
			return fmt.Errorf("Error al desinstalar la versión actual: %w", err)
		}
		fmt.Println("La versión actual ha sido desinstalada. Por favor, selecciona una nueva versión usando 'use'")
	}

	fmt.Printf("La versión %s ha sido desinstalada correctamente\n", version)
	return nil
}
