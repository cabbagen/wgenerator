package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
)

func HandleFeignRequest(uri, method string, params map[string]string, headers map[string]string) (interface{}, error) {
	client := resty.New().R().SetHeaders(headers)

	// GET 方法相关
	if method == resty.MethodGet || method == resty.MethodHead || method == resty.MethodDelete {
		client.SetQueryParams(params)
	}

	// POST 方法相关
	if (method == resty.MethodPost || method == resty.MethodPut) && headers["Content-Type"] == "application/json" {
		client.SetBody(params)
	}

	if (method == resty.MethodPost || method == resty.MethodPut) && headers["content-type"] == "application/x-www-form-urlencoded" {
		client.SetFormData(params)
	}

	// 发送请求
	repsonse, error := client.Execute(method, uri)

	if error != nil {
		return nil, errors.New("发送请求错误：" + error.Error())
	}

	if !repsonse.IsSuccess() {
		return nil, errors.New("本次请求响应失败")
	}

	result, resultJSON := repsonse.Body(), make(map[string]string)

	if strings.HasPrefix(repsonse.Header().Get("Content-Type"), "application/json") {
		error := json.Unmarshal(result, &resultJSON)
		return resultJSON, error
	}

	return result, nil
}

type RequestFile struct {
	file     []byte
	Field    string
	FileName string
}

func HandleFeignFileRequest(uri string, files []RequestFile, formData map[string]string, headers map[string]string) (interface{}, error) {
	client := resty.New().R().SetHeaders(headers)

	// 绑定文件
	for _, file := range files {
		client.SetFileReader(file.Field, file.FileName, bytes.NewReader(file.file))
	}

	// 发送请求
	repsonse, error := client.SetFormData(formData).Post(uri)

	if error != nil {
		return nil, errors.New("发送请求错误：" + error.Error())
	}

	if !repsonse.IsSuccess() {
		return nil, errors.New("本次请求响应失败")
	}

	result, resultJSON := repsonse.Body(), make(map[string]string)

	if strings.HasPrefix(repsonse.Header().Get("Content-Type"), "application/json") {
		error := json.Unmarshal(result, &resultJSON)
		return resultJSON, error
	}

	return result, nil
}
