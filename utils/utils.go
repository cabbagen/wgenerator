package utils

import (
	"encoding/json"
	"strings"
)

// 字符串输出
type Outputable interface {
	ToString() string
}

/**
 * 从 Slice 结构中，查找第一个符合条件的元素
 */
func FindBySlice[T any](list []T, fn func(value T, index int) bool) (result T, isExist bool) {
	for index, value := range list {
		if fn(value, index) {
			result, isExist = value, true
			return
		}
	}
	return
}

/**
 * 从 Slice 结构中，map 产生一个新的 Slice
 */
func MapBySlice[T any, K any](list []T, fn func(value T, index int) K) (results []K) {
	for index, value := range list {
		results = append(results, fn(value, index))
	}
	return
}

/**
 * 循环遍历 Slice 结构
 */
func ForEachBySlice[T any](list []T, fn func(value T, index int)) {
	for index, value := range list {
		fn(value, index)
	}
}

/**
 * 从 Slice 结构中，filter 产生一个新的 Slice
 */
func FilterBySlice[T any](list []T, fn func(value T, index int) bool) (results []T) {
	for index, value := range list {
		if fn(value, index) {
			results = append(results, value)
		}
	}
	return
}

/**
 * 获取 Slice 中的目标索引， -1 为不存在
 */
func IndexOfBySlice[T comparable](list []T, target T) (targetIndex int) {
	targetIndex = -1

	for index, value := range list {
		if value == target {
			targetIndex = index
			return
		}
	}
	return
}

/*
 * 检查元素是否存在
 **/
func IsExistBySlice[T comparable](list []T, target T) bool {
	return IndexOfBySlice(list, target) > -1
}

/**
 * 根据索引，查找 Slice 中的目标元素
 */
func AtBySlice[T any](list []T, index int) (result T, isExist bool) {
	if len(list) <= index || index < 0 {
		isExist = false
		return
	}
	result, isExist = list[index], true
	return
}

/**
 * 拼接 Slice
 */
func ConcatBySlice[T any](list []T, appendList []T) []T {
	return append(list, appendList...)
}

/**
 * 拼接 Slice 为字符串返回
 */
func JoinBySlice[T Outputable](list []T, separator string) string {
	willJoinStrings := MapBySlice[T, string](list, func(value T, _ int) string {
		return value.ToString()
	})
	return strings.Join(willJoinStrings, separator)
}

/**
 * struct 转 map 结构
 */
func StructToMap(src interface{}) (map[string]interface{}, error) {
	data, error := json.Marshal(src)

	if error != nil {
		return make(map[string]interface{}), error
	}

	result := make(map[string]interface{})

	if error := json.Unmarshal(data, &result); error != nil {
		return result, error
	}

	return result, nil
}
