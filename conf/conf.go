/**
 * 配置文件操作相关
 * ==========================================================
 * 配置文件规则: [debug|test|release].config.yaml
 */
package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

var workspaceDir string = "."

func init() {
	if workspaceDir != "." {
		return
	}
	if dir, err := os.Getwd(); err == nil {
		workspaceDir = dir
	}
}

func GetYamlFilePath() string {
	mode := gin.Mode()

	if mode == "" {
		mode = "debug"
	}
	return filepath.Join(workspaceDir, fmt.Sprintf("./%s.config.yaml", mode))
}

func ScanfBuildinYamlConfig() (config ApplicationConfig, error error) {
	datas, error := os.ReadFile(GetYamlFilePath())

	if error != nil {
		return config, error
	}

	if error := yaml.Unmarshal(datas, &config); error != nil {
		return config, error
	}

	return config, nil
}
