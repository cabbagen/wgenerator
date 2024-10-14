/**
 * 配置文件操作相关
 * ==========================================================
 * 配置文件规则: [debug|test|release|$custome].config.yaml
 */
package conf

import (
	"fmt"
	"os"
	"wgenerator/definitions"

	"gopkg.in/yaml.v3"
)

func GetYamlFilePath() string {
	mode := os.Getenv(definitions.WGeneratorENV)

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
