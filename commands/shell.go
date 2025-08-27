package commands

import (
	"fmt"
	"os"
	"os/exec"
	"polynode/shared"
	"runtime"
)

func OpenShell() error {
	// Obtener la versión actual de Node.js
	currentVersion := shared.GetCurrentVersion()
	if currentVersion == "" {
		return fmt.Errorf("no hay ninguna versión seleccionada. Utilice el comando 'use' para seleccionar una versión")
	}

	// Obtener la ruta del directorio actual de Node.js
	currentVersionPath := shared.GetCurrentVersionPath()
	
	// Verificar que el directorio existe
	if _, err := os.Stat(currentVersionPath); os.IsNotExist(err) {
		return fmt.Errorf("la versión actual no está instalada correctamente")
	}

	// Configurar las variables de entorno
	env := os.Environ()
	
	// Agregar el directorio de Node.js al PATH
	nodePath := currentVersionPath
	pathEnv := os.Getenv("PATH")
	
	// En Windows, el PATH se separa con punto y coma
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	
	// Agregar el directorio de Node.js al inicio del PATH
	newPath := nodePath + separator + pathEnv
	env = append(env, "PATH="+newPath)
	
	// Agregar variable de entorno para indicar que estamos en un shell de polynode
	env = append(env, "POLYNODE_SHELL=true")
	env = append(env, fmt.Sprintf("POLYNODE_VERSION=%s", currentVersion))

	// Determinar el shell por defecto según el sistema operativo
	var shellCmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		// En Windows, usar cmd.exe
		shellCmd = exec.Command("cmd.exe")
		shellCmd.Env = env
		shellCmd.Stdin = os.Stdin
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr
	} else {
		// En Unix/Linux, usar el shell por defecto del usuario
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		shellCmd = exec.Command(shell)
		shellCmd.Env = env
		shellCmd.Stdin = os.Stdin
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr
	}

	fmt.Printf("Abriendo shell con Node.js v%s...\n", currentVersion)
	fmt.Printf("Directorio de Node.js: %s\n", currentVersionPath)
	fmt.Println("Presiona Ctrl+C para salir del shell")

	// Ejecutar el shell
	return shellCmd.Run()
}
