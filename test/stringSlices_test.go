package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestStringSlices(t *testing.T) {
	//str := "[1,2]"
	//fmt.Println(str, str[0])
	//strings.TrimLeft(strings.TrimRight(str[0], "]"), "[")
	//s := str[1 : len(str)-1]
	//slices := strings.Split(s, ",")
	//for _, slice := range slices {
	//	fmt.Println(slice, reflect.TypeOf(slice))
	//}

	// 给定的字符串
	slice := []string{
		`{"input":"1 2\n","output":"3"},{"input":"1 2\n","output":"3"}`,
	}
	str := slice[0]
	fmt.Println(reflect.TypeOf(slice[0]))

	//// 使用正则表达式匹配 JSON 对象
	//re := regexp.MustCompile(`{[^}]+}`)
	//jsonStrings := re.FindAllString(str, -1)
	//
	//// 遍历每个 JSON 对象
	//for _, jsonString := range jsonStrings {
	//	// 解析 JSON 字符串到结构体
	//	var data map[string]interface{}
	//	err := json.Unmarshal([]byte(jsonString), &data)
	//	if err != nil {
	//		fmt.Println("Error parsing JSON:", err)
	//		continue
	//	}
	//	fmt.Println(data)
	//}
	// 通过逗号分割字符串
	parts := strings.Split(str, "},{")

	// 修正第一个和最后一个片段，以确保他们是单独的 JSON 对象
	parts[0] = strings.TrimPrefix(parts[0], "{")
	parts[len(parts)-1] = strings.TrimSuffix(parts[len(parts)-1], "}")

	// 打印两个 JSON 对象
	testcaselist := make([]string, 0)
	for _, part := range parts {
		jsonStr := "{" + part + "}"
		jsonStr = "[" + jsonStr + "]"
		fmt.Println(jsonStr)
		testcaselist = append(testcaselist, jsonStr)
	}
	fmt.Println(testcaselist)

}
