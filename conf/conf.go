/**
 * 配置文件操作相关
 * ==========================================================
 * 配置文件规则: [debug|test|release].config.yaml
 */
package conf

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func GetYamlFilePath() string {
	mode := gin.Mode()

	if mode == "" {
		mode = "debug"
	}
	return fmt.Sprintf("./%s.config.yaml", mode)
}

func ScanfYamlConfig(confPath string) (config map[string]map[string]interface{}, error error) {
	datas, error := os.ReadFile(confPath)

	if error != nil {
		return config, error
	}

	if error := yaml.Unmarshal(datas, &config); error != nil {
		return config, error
	}

	return config, nil
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
