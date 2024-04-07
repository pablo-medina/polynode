package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type version struct {
	Major int
	Minor int
	Patch int
}

const nodeURLTemplate = "https://nodejs.org/dist/v%s/node-v%s-%s-x64.zip"

var installPath string

func init() {
	installPath = os.Getenv("POLYNODE_PATH")
	if installPath == "" {
		installPath = "C:\\polynode"
	}

	// Verificar si la carpeta de instalación ya existe
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		// La carpeta no existe, crearla
		if err := os.MkdirAll(installPath, 0755); err != nil {
			fmt.Printf("Error al crear la carpeta de instalación: %v\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: polynode <init|install|use|list|version> [version]")
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		err := initPolyNode()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Configuración inicial de Polynode completada")
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Uso: polynode install <version>")
			return
		}
		version := os.Args[2]
		err := installNode(version)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Node v%s instalado en %s\n", version, installPath)

	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Uso: polynode use <version>")
			return
		}
		version := os.Args[2]
		err := useNode(version)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Se cambió la versión a v%s\n", version)

	case "list":
		listInstalledVersionsCommand()

	case "version":
		versionCommand()

	default:
		fmt.Println("Uso: polynode <init|install|use|list|version> [version]")
	}
}

func initPolyNode() error {
	setenvPath := filepath.Join(installPath, "setenv.cmd")

	// Crear el archivo setenv.cmd
	f, err := os.Create(setenvPath)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo setenv.cmd: %v", err)
	}
	defer f.Close()

	// Escribir el contenido del archivo setenv.cmd
	_, err = f.WriteString(fmt.Sprintf("@ECHO OFF\nset PATH=%s\\current;%%PATH%%\n", installPath))
	if err != nil {
		return fmt.Errorf("Error al escribir en el archivo setenv.cmd: %v", err)
	}

	return nil
}

func installNode(version string) error {
	zipURL := fmt.Sprintf(nodeURLTemplate, version, version, getOS())
	resp, err := http.Get(zipURL)
	if err != nil {
		return fmt.Errorf("Error al obtener el archivo ZIP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("No se pudo encontrar la versión especificada: %s", version)
	}

	// Crear el directorio de instalación si no existe
	err = os.MkdirAll(installPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error al crear el directorio de instalación: %v", err)
	}

	zipFileName := filepath.Join(installPath, fmt.Sprintf("node-v%s-%s-x64.zip", version, getOS()))

	// Guardar el archivo ZIP
	outFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo ZIP: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("Error al guardar el archivo ZIP: %v", err)
	}

	// Extraer el archivo ZIP
	err = unzip(zipFileName, installPath)
	if err != nil {
		return fmt.Errorf("Error al extraer el archivo ZIP: %v", err)
	}

	// Cerrar el archivo ZIP antes de intentar eliminarlo
	outFile.Close()

	// Crear un archivo version.info
	versionInfoPath := filepath.Join(installPath, fmt.Sprintf("node-v%s-win-x64", version), "version.info")
	err = os.WriteFile(versionInfoPath, []byte(version), 0644)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo version.info: %v", err)
	}

	// Eliminar el archivo ZIP después de extraerlo
	err = os.Remove(zipFileName)
	if err != nil {
		return fmt.Errorf("Error al eliminar el archivo ZIP: %v", err)
	}

	return nil
}

func useNode(version string) error {
	versionPath := filepath.Join(installPath, fmt.Sprintf("node-v%s-win-x64", version))
	currentPath := filepath.Join(installPath, "current")

	// Verificar si la carpeta de la versión de Node existe
	_, err := os.Stat(versionPath)
	if err != nil {
		return fmt.Errorf("La versión especificada de Node no está instalada: %s", version)
	}

	// Mover la versión anterior si es necesario
	err = movePrevious(version, currentPath)
	if err != nil {
		return fmt.Errorf("Error al mover la versión anterior: %v", err)
	}

	// Eliminar el directorio actual si ya existe
	if _, err := os.Stat(currentPath); !os.IsNotExist(err) {
		os.RemoveAll(currentPath)
	}

	// Copiar el contenido completo del directorio de la versión a current
	err = copyDir(versionPath, currentPath)
	if err != nil {
		return fmt.Errorf("Error al copiar los archivos: %v", err)
	}

	return nil
}

func listInstalledVersionsCommand() {
	versions, err := listInstalledVersions()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(versions) == 0 {
		fmt.Println("No se encontraron versiones instaladas.")
		return
	}

	fmt.Println("Versiones instaladas:")
	for _, version := range versions {
		fmt.Println("  -", version)
	}
}

func versionCommand() {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Versión seleccionada:", currentVersion)
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

func getCurrentVersion() (string, error) {
	versionFilePath := filepath.Join(installPath, "current", "version.info")

	// Leer el contenido del archivo version.info
	versionBytes, err := os.ReadFile(versionFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("No hay una versión seleccionada")
		}
		return "", fmt.Errorf("Error al leer el archivo version.info: %v", err)
	}

	return strings.TrimSpace(string(versionBytes)), nil
}

func movePrevious(version string, currentPath string) error {
	// Leer el archivo version.info dentro de current
	currentVersionFile := filepath.Join(currentPath, "version.info")
	currentVersionBytes, err := os.ReadFile(currentVersionFile)
	if err != nil {
		if os.IsNotExist(err) {
			// No hay una versión anterior, no es necesario hacer nada
			return nil
		}
		return fmt.Errorf("Error al leer el archivo version.info en current: %v", err)
	}
	currentVersion := strings.TrimSpace(string(currentVersionBytes))

	// Verificar si la versión actual es diferente
	if currentVersion != version {
		// Eliminar la carpeta de la versión anterior si existe
		previousPath := filepath.Join(installPath, fmt.Sprintf("node-v%s-win-x64", currentVersion))
		if _, err := os.Stat(previousPath); !os.IsNotExist(err) {
			os.RemoveAll(previousPath)
		}

		// Renombrar current con el nombre de la versión anterior
		err := os.Rename(currentPath, previousPath)
		if err != nil {
			return fmt.Errorf("Error al renombrar current: %v", err)
		}
	}

	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getOS() string {
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		return "win"
	}
	return "linux"
}

func parseVersion(versionStr string) (version, error) {
	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		return version{}, fmt.Errorf("Formato de versión inválido: %s", versionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return version{}, fmt.Errorf("Error al convertir la parte mayor: %v", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return version{}, fmt.Errorf("Error al convertir la parte menor: %v", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return version{}, fmt.Errorf("Error al convertir la parte de revisión: %v", err)
	}

	return version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func compareVersions(v1, v2 version) int {
	if v1.Major != v2.Major {
		return v1.Major - v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor - v2.Minor
	}
	return v1.Patch - v2.Patch
}

func listInstalledVersions() ([]string, error) {
	// Obtener una lista de todos los directorios dentro de la carpeta de instalación
	files, err := os.ReadDir(installPath)
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
