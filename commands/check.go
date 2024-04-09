package commands

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"polynode/shared"
	"strings"
)

func CheckInstallation() error {
	currentPath := shared.GetCurrentVersionPath()

	// Intentar encontrar la ubicación de "where" en el PATH actual
	whereCmd := exec.Command("where", "where")
	whereCmd.Stdout = nil
	whereCmd.Stderr = nil
	if err := whereCmd.Run(); err != nil {
		// Si no se encuentra "where" en el PATH, intentar desde la ruta "c:\windows\system32"
		whereCmd = exec.Command(filepath.Join("c:\\windows\\system32", "where"), "where")
	}

	// Ejecutar "where node" para obtener la lista de ubicaciones del ejecutable de Node.js
	cmd := exec.Command(whereCmd.Path, "node")
	output, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error al ejecutar 'where node': %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error al iniciar el comando 'where node': %w", err)
	}

	// Leer la salida del comando línea por línea
	scanner := bufio.NewScanner(output)
	var nodePaths []string
	for scanner.Scan() {
		nodePaths = append(nodePaths, strings.TrimSpace(scanner.Text()))
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error al esperar el comando 'where node': %w", err)
	}

	// Verificar la lista de ubicaciones del ejecutable de Node.js
	switch len(nodePaths) {
	case 0:
		// Si no hay elementos en la lista, indicar cómo actualizar el PATH
		fmt.Println("No se encontraron ubicaciones para el ejecutable de Node.js.")
		fmt.Println("Por favor, actualice el PATH del sistema con el siguiente comando:")
		fmt.Printf("SET PATH=%s;%%PATH%%\n", currentPath)
	case 1:
		// Si hay una sola ubicación, verificar si coincide con la ubicación esperada
		if nodePaths[0] != filepath.Join(currentPath, "node.exe") {
			fmt.Println("La ubicación del ejecutable de Node.js no coincide con la ubicación esperada.")
			fmt.Println("Por favor, actualice el PATH del sistema con el siguiente comando:")
			fmt.Printf("SET PATH=%s;%%PATH%%\n", currentPath)
		} else {
			fmt.Println("La ubicación del ejecutable de Node.js coincide con la ubicación esperada.")
		}
	default:
		// Si hay más de una ubicación, verificar si la primera coincide con la ubicación esperada
		if nodePaths[0] != filepath.Join(currentPath, "node.exe") {
			fmt.Println("La primera ubicación del ejecutable de Node.js no coincide con la ubicación esperada.")
			fmt.Println("Por favor, actualice el PATH del sistema con el siguiente comando:")
			fmt.Printf("SET PATH=%s;%%PATH%%\n", currentPath)
		} else {
			fmt.Println("La primera ubicación del ejecutable de Node.js coincide con la ubicación esperada.")
			fmt.Println("Se encontraron múltiples ubicaciones del ejecutable de Node. Se recomienda eliminar las ubicaciones adicionales del PATH correspondientes a otras instalaciones de Node para evitar conflictos.")
		}
	}

	return nil
}
