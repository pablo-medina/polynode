package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envAppPath             = "POLYNODE_PATH"
	defaultInstallPath     = "C:\\polynode"
	currentVersionPathName = "current"
	repoPathName           = "repository"
	nodeURLTemplate        = "https://nodejs.org/dist/v%s/node-v%s-%s-x64.zip"
)

var (
	installPath        string
	currentVersionPath string
	repoPath           string
)

type Version struct {
	Major int
	Minor int
	Patch int
}

type ProxyConfig struct {
	HTTPProxy  string `json:"http_proxy,omitempty"`
	HTTPSProxy string `json:"https_proxy,omitempty"`
}

func init() {
	/*
	   Initialize InstallPath with the environment variable "POLYNODE_PATH".
	   If it doesn't exist, use the DefaultInstallPath constant instead.
	*/
	installPath = os.Getenv(envAppPath)
	if installPath == "" {
		installPath = defaultInstallPath
	}

	currentVersionPath = filepath.Join(installPath, currentVersionPathName)
	repoPath = filepath.Join(installPath, repoPathName)
}

func GetInstallPath() string {
	return installPath
}

func GetCurrentVersionPath() string {
	return currentVersionPath
}

func GetRepoPath() string {
	return repoPath
}

func getOS() string {
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		return "win"
	}
	return "linux"
}

func GetNodeVersionURL(version string) string {
	return fmt.Sprintf(nodeURLTemplate, version, version, getOS())
}

func GetCurrentVersion() string {
	// Verificar si el directorio current existe
	_, err := os.Stat(currentVersionPath)
	if err != nil && os.IsNotExist(err) {
		return ""
	}

	// Obtener la ruta del archivo version.info
	versionFile := filepath.Join(currentVersionPath, "version.info")

	// Verificar si el archivo version.info existe
	_, err = os.Stat(versionFile)
	if err != nil && os.IsNotExist(err) {
		return ""
	}

	// Leer el contenido del archivo version.info
	versionBytes, err := os.ReadFile(versionFile)
	if err != nil {
		fmt.Println("Error al leer la versión actual:", err)
		return ""
	}

	return strings.TrimSpace(string(versionBytes))
}
