package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"polynode/shared"
)

func UseNodeVersion(version string) error {
	versionPath := filepath.Join(shared.GetRepoPath(), fmt.Sprintf("node-v%s-win-x64", version))

	// Verificar si la carpeta de la versión de Node existe
	_, err := os.Stat(versionPath)
	if err != nil {
		return fmt.Errorf("La versión especificada de Node no está instalada: %s", version)
	}

	// Leer la versión actualmente seleccionada
	currentVersion := shared.GetCurrentVersion()

	// Comprobar si la versión solicitada es la misma que la actual
	if currentVersion == version {
		return fmt.Errorf("La versión %s ya está seleccionada, no es necesario cambiar de versión", version)
	} else {
		if currentVersion != "" {
			// Mover la versión anterior si es necesario
			err = movePrevious(version)
			if err != nil {
				return fmt.Errorf("Error al mover la versión anterior: %v", err)
			}
		}
	}

	// Eliminar el directorio actual si ya existe
	if _, err := os.Stat(shared.GetCurrentVersionPath()); !os.IsNotExist(err) {
		os.RemoveAll(shared.GetCurrentVersionPath())
	}

	// Copiar el contenido completo del directorio de la versión a current
	err = copyDir(versionPath, shared.GetCurrentVersionPath())
	if err != nil {
		return fmt.Errorf("Error al copiar los archivos: %v", err)
	}

	return nil
}

func movePrevious(version string) error {
	currentVersion := shared.GetCurrentVersion()

	// Verificar si la versión actual es diferente
	if version != "" && currentVersion != version {
		// Eliminar la carpeta de la versión anterior si existe
		previousPath := filepath.Join(shared.GetRepoPath(), fmt.Sprintf("node-v%s-win-x64", currentVersion))
		if _, err := os.Stat(previousPath); !os.IsNotExist(err) {
			os.RemoveAll(previousPath)
		}

		// Renombrar current con el nombre de la versión anterior
		err := os.Rename(shared.GetCurrentVersionPath(), previousPath)
		if err != nil {
			return fmt.Errorf("Error al renombrar current: %v", err)
		}
	}

	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcFilePath := filepath.Join(src, file.Name())
		dstFilePath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			if err := copyDir(srcFilePath, dstFilePath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcFilePath, dstFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
