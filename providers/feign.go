package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
)

func handleRestyResponseAdaptor(response *resty.Response) (interface{}, error) {
	if !response.IsSuccess() {
		return nil, errors.New("本次请求响应失败: " + response.Status())
	}

	result, resultJSON := response.Body(), make(map[string]interface{})

	if strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
		error := json.Unmarshal(result, &resultJSON)
		return resultJSON, error
	}

	return result, nil
}

func HandleFeignProxyRequest(uriWithParams, method string, body []byte, headers map[string]string) (interface{}, error) {
	client := resty.New().R().SetHeaders(headers)

	// POST 方法相关
	if method == resty.MethodPost || method == resty.MethodPut {
		client.SetBody(body)
	}

	// 发送请求
	response, error := client.Execute(method, uriWithParams)

	if error != nil {
		return nil, errors.New("发送请求错误：" + error.Error())
	}

	return handleRestyResponseAdaptor(response)
}

func HandleFeignRequest(uri, method string, params interface{}, headers map[string]string) (interface{}, error) {
	client := resty.New().R().SetHeaders(headers)

	// GET 方法相关
	if method == resty.MethodGet || method == resty.MethodHead || method == resty.MethodDelete {
		client.SetQueryParams(params.(map[string]string))
	}

	// POST 方法相关
	if (method == resty.MethodPost || method == resty.MethodPut) && headers["Content-Type"] == "application/json" {
		client.SetBody(params)
	}

	if (method == resty.MethodPost || method == resty.MethodPut) && headers["Content-Type"] == "application/x-www-form-urlencoded" {
		client.SetFormData(params.(map[string]string))
	}

	// 发送请求
	response, error := client.Execute(method, uri)

	if error != nil {
		return nil, errors.New("发送请求错误：" + error.Error())
	}

	return handleRestyResponseAdaptor(response)
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
	response, error := client.SetFormData(formData).Post(uri)

	if error != nil {
		return nil, errors.New("发送请求错误：" + error.Error())
	}

	return handleRestyResponseAdaptor(response)
}

func HandleFeignPutFileRequest(uri string, thunk []byte, headers map[string]string) (bool, error) {
	if _, isExist := headers["Content-Type"]; !isExist {
		headers["Content-Type"] = "application/octet-stream"
	}

	response, error := resty.New().R().SetHeaders(headers).SetBody(thunk).Put(uri)

	if error != nil {
		return false, errors.New("发送请求错误：" + error.Error())
	}

	if !response.IsSuccess() {
		return false, errors.New("本次请求响应失败: " + response.Status())
	}

	return true, nil
}
