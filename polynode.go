package main

import (
	"fmt"
	"os"
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
		commands.ShowHelp()
		return
	}

	command := os.Args[1]

	switch command {

	case "check":
		err := commands.CheckInstallation()
		if err != nil {
			fmt.Println(err)
			return
		}
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Uso: poly install <version>")
			return
		}
		version := os.Args[2]
		err := commands.InstallVersion(version)
		if err != nil {
			fmt.Println(err)
			return
		}

	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Uso: poly use <version>")
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
			fmt.Println("Uso: poly uninstall <version>")
			return
		}
		version := os.Args[2]
		if err := commands.UninstallNodeVersion(version); err != nil {
			fmt.Println("Error al desinstalar la versión:", err)
			return
		}

	case "proxy":
		if len(os.Args) < 3 {
			fmt.Println("Uso: poly proxy <http_proxy>")
			return
		}

		httpProxy := os.Args[2]

		if err := commands.SetProxyURL(httpProxy); err != nil {
			fmt.Println("Error al configurar el proxy:", err)
			return
		}

	case "backup":
		if err := commands.BackupInstallation(); err != nil {
			fmt.Println(err)
			return
		}

	case "shell":
		if err := commands.OpenShell(); err != nil {
			fmt.Println("Error al abrir el shell:", err)
			return
		}

	case "help":
		commands.ShowHelp()
		return

	default:
		commands.ShowHelp()
	}
}
