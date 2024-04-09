package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"polynode/shared"
	"time"
)

func BackupInstallation() error {
	installPath := shared.GetInstallPath()

	// Generar el nombre del directorio de backup con un timestamp
	timestamp := time.Now().Format("20060102T150405")
	backupFilename := fmt.Sprintf("%s-polynode", timestamp)
	zipFileName := filepath.Join(installPath, fmt.Sprintf("%s.zip", backupFilename))

	// Crear el archivo zip
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo zip: %w", err)
	}
	defer zipFile.Close()

	// Crear un escritor de zip
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Función para recorrer los directorios y agregar archivos al zip
	addFileToZip := func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.Mode().IsRegular() {
			return nil // Ignorar directorios y otros tipos de archivos
		}

		// Abrir el archivo
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Obtener el nombre relativo del archivo
		relativePath, err := filepath.Rel(installPath, filePath)
		if err != nil {
			return err
		}

		// Crear un nuevo archivo en el zip
		zipEntry, err := zipWriter.Create(filepath.Join("polynode", relativePath))
		if err != nil {
			return err
		}

		// Copiar el contenido del archivo al zip
		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	}

	// Recorrer el directorio "current" y agregar archivos al zip
	if err := filepath.Walk(filepath.Join(installPath, "current"), addFileToZip); err != nil {
		return fmt.Errorf("error al agregar archivos de 'current' al zip: %w", err)
	}

	// Recorrer el directorio "repository" y agregar archivos al zip
	if err := filepath.Walk(filepath.Join(installPath, "repository"), addFileToZip); err != nil {
		return fmt.Errorf("error al agregar archivos de 'repository' al zip: %w", err)
	}

	fmt.Printf("Se creó una copia de seguridad en: %s\n", zipFileName)
	return nil
}
