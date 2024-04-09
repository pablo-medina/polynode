package commands

import (
	"fmt"
	"os"
	"polynode/shared"
	"sort"
	"strconv"
	"strings"
)

func ExecuteList() {
	versions, err := listInstalledVersions()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(versions) == 0 {
		fmt.Println("No se encontraron versiones instaladas.")
		return
	}

	currentVersion := shared.GetCurrentVersionPath()

	fmt.Println("Versiones instaladas:")
	for _, version := range versions {
		versionLine := ""
		if version == currentVersion {
			versionLine = fmt.Sprintf(" - [%s] <- ACTUAL", version)
		} else {
			versionLine = fmt.Sprintf(" - %s", version)
		}
		fmt.Println(versionLine)
	}
}

func listInstalledVersions() ([]string, error) {
	// Obtener una lista de todos los directorios dentro de la carpeta de instalación
	files, err := os.ReadDir(shared.GetRepoPath())
	if err != nil {
		return nil, fmt.Errorf("Error al leer la carpeta de instalación: %v", err)
	}

	var versions []string
	for _, file := range files {
		if file.IsDir() {
			// Excluir la carpeta 'current' de la lista
			if file.Name() != "current" {
				// Verificar si el nombre del directorio corresponde a una versión de Node
				if strings.HasPrefix(file.Name(), "node-v") {
					// Obtener solo el número de versión, sin el sufijo 'win-x64'
					version := strings.TrimPrefix(file.Name(), "node-v")
					version = strings.TrimSuffix(version, "-win-x64")
					versions = append(versions, version)
				}
			}
		}
	}

	// Convertir las versiones a una estructura personalizada y ordenarlas
	var sortedVersions []string
	for _, ver := range versions {
		_, err := parseVersion(ver)
		if err != nil {
			fmt.Printf("Error al analizar la versión %s: %v\n", ver, err)
			continue
		}
		sortedVersions = append(sortedVersions, ver)
	}
	sort.Slice(sortedVersions, func(i, j int) bool {
		v1, _ := parseVersion(sortedVersions[i])
		v2, _ := parseVersion(sortedVersions[j])
		return compareVersions(v1, v2) < 0
	})

	return sortedVersions, nil
}

func compareVersions(v1, v2 shared.Version) int {
	if v1.Major != v2.Major {
		return v1.Major - v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor - v2.Minor
	}
	return v1.Patch - v2.Patch
}

func parseVersion(versionStr string) (shared.Version, error) {
	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		return shared.Version{}, fmt.Errorf("Formato de versión inválido: %s", versionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return shared.Version{}, fmt.Errorf("Error al convertir la parte mayor: %v", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return shared.Version{}, fmt.Errorf("Error al convertir la parte menor: %v", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return shared.Version{}, fmt.Errorf("Error al convertir la parte de revisión: %v", err)
	}

	return shared.Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}
