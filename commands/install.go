package commands

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"polynode/shared"
	"strings"
)

var httpClient *http.Client

func getOS() string {
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		return "win"
	}
	return "linux"
}

func InstallVersion(version string) error {
	client := buildHttpClient()
	if client == nil {
		return fmt.Errorf("No se pudo procesar la configuración del proxy")
	}

	zipURL := shared.GetNodeVersionURL(version)
	fmt.Printf("Descargando archivo %s...\n", zipURL)

	req, err := http.NewRequest("GET", zipURL, nil)
	if err != nil {
		return fmt.Errorf("Error al crear la solicitud HTTP:", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error al obtener el archivo ZIP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("No se pudo encontrar la versión especificada: %s", version)
	}

	// Crear el directorio de instalación si no existe
	err = os.MkdirAll(shared.GetRepoPath(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error al crear el directorio de instalación: %v", err)
	}

	zipFileName := filepath.Join(shared.GetRepoPath(), fmt.Sprintf("node-v%s-%s-x64.zip", version, getOS()))

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
	err = unzip(zipFileName, shared.GetRepoPath())
	if err != nil {
		return fmt.Errorf("Error al extraer el archivo ZIP: %v", err)
	}

	// Cerrar el archivo ZIP antes de intentar eliminarlo
	outFile.Close()

	// Crear un archivo version.info
	versionInfoPath := filepath.Join(shared.GetRepoPath(), fmt.Sprintf("node-v%s-win-x64", version), "version.info")
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

func buildHttpClient() *http.Client {
	proxyConfigFile := filepath.Join(shared.GetInstallPath(), "proxy.json")

	// Leer la configuración del proxy desde el archivo proxy.json si existe
	proxyConfig := shared.ProxyConfig{}
	if _, err := os.Stat(proxyConfigFile); err == nil {
		file, err := os.Open(proxyConfigFile)
		if err != nil {
			fmt.Println("Error al abrir el archivo proxy.json:", err)
			return nil
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&proxyConfig); err != nil {
			fmt.Println("Error al decodificar el archivo proxy.json:", err)
			return nil
		}
	}

	// Configurar el cliente HTTP con el proxy si está definido
	if proxyConfig.HTTPProxy != "" {
		fmt.Println("Configuración de proxy detectada")
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyConfig.HTTPProxy)
		}
		transport := &http.Transport{Proxy: proxy}
		httpClient = &http.Client{Transport: transport}
	} else {
		httpClient = http.DefaultClient
	}

	return httpClient
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
