package main

import (
	"fmt"
	"os"
	"path/filepath"
	"polynode/commands"
	"polynode/shared"
)

type ProxyConfig struct {
	HTTPProxy  string `json:"http_proxy,omitempty"`
	HTTPSProxy string `json:"https_proxy,omitempty"`
}

func init() {
	// Verificar si la carpeta de instalación ya existe
	if _, err := os.Stat(shared.GetInstallPath()); os.IsNotExist(err) {
		// La carpeta no existe, crearla
		if err := os.MkdirAll(shared.GetInstallPath(), 0755); err != nil {
			fmt.Printf("Error al crear la carpeta de instalación: %v\n", err)
			os.Exit(1)
		}
	}

	// Verificar si la carpeta de repositorio ya existe
	if _, err := os.Stat(shared.GetRepoPath()); os.IsNotExist(err) {
		// La carpeta no existe, crearla
		if err := os.MkdirAll(shared.GetRepoPath(), 0755); err != nil {
			fmt.Printf("Error al crear la carpeta de instalación: %v\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: polynode <init|install|use|list|version|uninstall|proxy> [version]")
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
		err := commands.InstallVersion(version)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Node v%s instalado en %s\n", version, shared.GetInstallPath())

	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Uso: polynode use <version>")
			return
		}
		version := os.Args[2]
		err := commands.UseNodeVersion(version)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Se cambió la versión a v%s\n", version)

	case "list":
		commands.ExecuteList()

	case "version":
		commands.ShowCurrentNodeVersion()

	case "uninstall":
		if len(os.Args) < 3 {
			fmt.Println("Uso: polynode uninstall <version>")
			return
		}
		version := os.Args[2]
		if err := commands.UninstallNodeVersion(version); err != nil {
			fmt.Println("Error al desinstalar la versión:", err)
			return
		}

	case "proxy":
		if len(os.Args) < 3 {
			fmt.Println("Uso: polynode proxy <http_proxy>")
			return
		}

		httpProxy := os.Args[2]

		if err := commands.SetProxyURL(httpProxy); err != nil {
			fmt.Println("Error al configurar el proxy:", err)
			return
		}

	default:
		fmt.Println("Uso: polynode <init|install|use|list|version|uninstall|proxy> [version]")
	}
}

func initPolyNode() error {
	setenvPath := filepath.Join(shared.GetInstallPath(), "setenv.cmd")

	// Crear el archivo setenv.cmd
	f, err := os.Create(setenvPath)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo setenv.cmd: %v", err)
	}
	defer f.Close()

	// Escribir el contenido del archivo setenv.cmd
	_, err = f.WriteString(fmt.Sprintf("@ECHO OFF\nset PATH=%s\\current;%%PATH%%\n", shared.GetInstallPath()))
	if err != nil {
		return fmt.Errorf("Error al escribir en el archivo setenv.cmd: %v", err)
	}

	return nil
}
