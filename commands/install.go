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

type IndexEntry struct {
	Version string      `json:"version"`
	Lts     interface{} `json:"lts"`
}

// ProgressReader es un wrapper para io.Reader que muestra progreso
type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	FileName string
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Current += int64(n)
	pr.showProgress()
	return n, err
}

func (pr *ProgressReader) showProgress() {
	if pr.Total <= 0 {
		return
	}

	percentage := float64(pr.Current) / float64(pr.Total) * 100
	barLength := 30
	filledLength := int(float64(barLength) * percentage / 100)

	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar += "="
		} else if i == filledLength {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"

	// Calcular tamaños en MB
	currentMB := float64(pr.Current) / (1024 * 1024)
	totalMB := float64(pr.Total) / (1024 * 1024)

	// Limpiar línea anterior y mostrar progreso
	fmt.Printf("\r%s %s %.1f%% (%.1f/%.1f MB)", pr.FileName, bar, percentage, currentMB, totalMB)

	if pr.Current >= pr.Total {
		fmt.Println() // Nueva línea al completar
	}
}

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

	parsedVersion := version

	if version == "lts" {
		ltsVersion := getLatestLTSURL(client)
		if ltsVersion == "" {
			return fmt.Errorf("No se pudo obtener la versión LTS")
		}
		parsedVersion = ltsVersion
		fmt.Printf("Versión LTS: %s\n", parsedVersion)
	}

	zipURL := shared.GetNodeVersionURL(parsedVersion)
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
		return fmt.Errorf("No se pudo encontrar la versión especificada: %s", parsedVersion)
	}

	// Crear el directorio de instalación si no existe
	err = os.MkdirAll(shared.GetRepoPath(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error al crear el directorio de instalación: %v", err)
	}

	zipFileName := filepath.Join(shared.GetRepoPath(), fmt.Sprintf("node-v%s-%s-x64.zip", parsedVersion, getOS()))

	// Guardar el archivo ZIP con barra de progreso
	outFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo ZIP: %v", err)
	}
	defer outFile.Close()

	// Crear el ProgressReader para mostrar el progreso
	progressReader := &ProgressReader{
		Reader:   resp.Body,
		Total:    resp.ContentLength,
		FileName: filepath.Base(zipFileName),
	}

	_, err = io.Copy(outFile, progressReader)
	if err != nil {
		return fmt.Errorf("Error al guardar el archivo ZIP: %v", err)
	}

	fmt.Println("Extrayendo archivos...")

	// Extraer el archivo ZIP
	err = unzip(zipFileName, shared.GetRepoPath())
	if err != nil {
		return fmt.Errorf("Error al extraer el archivo ZIP: %v", err)
	}

	// Cerrar el archivo ZIP antes de intentar eliminarlo
	outFile.Close()

	// Eliminar el archivo ZIP después de extraerlo
	err = os.Remove(zipFileName)
	if err != nil {
		return fmt.Errorf("Error al eliminar el archivo ZIP: %v", err)
	}

	fmt.Printf("Node v%s instalado en %s\n", parsedVersion, shared.GetInstallPath())

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

func getLatestLTSURL(client *http.Client) string {
	baseURL := shared.GetNodeRepositoryBaseURL()
	jsonDataURL := baseURL + "index.json"

	req, err := http.NewRequest("GET", jsonDataURL, nil)
	if err != nil {
		fmt.Printf("Error al crear la solicitud HTTP: %s\n", err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error al realizar la solicitud HTTP para obtener información sobre las versiones disponibles: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error al obtener información sobre las versiones dipsonibles.")
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer el contenido del archivo index.json:", err)
		return ""
	}

	var versions []IndexEntry
	if err := json.Unmarshal(body, &versions); err != nil {
		fmt.Println("Error al decodificar el contenido del archivo index.json:", err)
		return ""
	}

	/*
	 * El campo lts puede ser un boolean con el valor false o un string con un código de versión lts
	 * Lo que hacemos es, recorrer el json, y devolver el primer elemento cuyo valor de "lts" sea un string.
	 */
	for _, entry := range versions {
		ltsValue := "false"
		switch lts := entry.Lts.(type) {
		case bool:
			if lts {
				ltsValue = "true"
			}
		case string:
			ltsValue = lts
		}
		if ltsValue != "false" {
			return shared.NormalizeVersion(entry.Version)
		}
	}

	return ""
}
