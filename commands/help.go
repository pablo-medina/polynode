package commands

import (
	"fmt"
)

func ShowHelp() {
	fmt.Println("Uso: poly <comando> [parámetros]")
	fmt.Println("")
	fmt.Println("Comandos:")
	fmt.Println("---------")
	fmt.Println(" init                   Inicializar polynode")
	fmt.Println(" install <version>      Instalar versión de node especificada en el repositorio local")
	fmt.Println(" use <version>          Usar versión de node previamente instalada")
	fmt.Println(" list                   Lista versiones de node instaladas en el repositorio local")
	fmt.Println(" version                Muestra la versión de Node seleccionada")
	fmt.Println(" uninstall <version>    Eliminar del repositorio local la versión de node especificada")
	fmt.Println(" proxy <url>            Utilizar la url de proxy indicada para la descarga de versiones de Node")
	fmt.Println(" help                   Mostrar esta ayuda")
	fmt.Println()
}
