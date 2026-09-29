/**
 * 配置文件操作相关
 * ==========================================================
 * 配置文件规则: [debug|test|release].config.yaml
 */
package conf

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func GetYamlFilePath() string {
	mode := gin.Mode()

	if mode == "" {
		mode = "debug"
	}

	_, file, _, _ := runtime.Caller(0)

	root := filepath.Dir(filepath.Dir(file))

	return filepath.Join(root, mode+".config.yaml")
}

func ScanfBuildinYamlConfig() (config map[string]map[string]interface{}, error error) {
	datas, error := os.ReadFile(GetYamlFilePath())

	if error != nil {
		return config, error
	}

	if error := yaml.Unmarshal(datas, &config); error != nil {
		return config, error
	}

	return config, nil
}
