package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 判断文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// 该字符串是否为数字
func isNumber(value interface{}) (int, bool) {
	s := fmt.Sprintf("%v", value)
	reg := regexp.MustCompile(`^-?[1-9]\d*$|0`)
	isNumber := reg.MatchString(s)
	if isNumber {
		num, _ := strconv.Atoi(s)
		return num, true
	}
	return 0, false
}

// 字符串是否为浮点类型
func isFloat64(value interface{}) (float64, bool) {
	s := fmt.Sprintf("%v", value)
	reg := regexp.MustCompile(`^-?([1-9]\d*\.?\d*|0\.\d*[1-9]\d*|0?\.0+|0)$`)
	isFloat := reg.MatchString(s)
	if isFloat {
		f, _ := strconv.ParseFloat(s, 64)
		return f, true
	}
	return 0, false
}

// 将形如[a.b[0].c[1]]拆解成列表["a", "b", "0", "c", "1"]
func splitKey(key string) []string {
	key = strings.ReplaceAll(key, "[", ".")
	key = strings.ReplaceAll(key, "]", ".")
	key = strings.TrimSuffix(key, ".")
	list := strings.Split(key, ".")
	return list
}

// splitKey的逆过程
func composeKey(l []string) string {
	var key string
	for index, item := range l {
		if index == 0 {
			key = item
			continue
		}
		_ , isnumber := isNumber(item)
		if isnumber {
			key += "[" + item + "]"
		} else {
			key += "." + item
		}
	}
	return key
}