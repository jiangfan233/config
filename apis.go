package config

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// 加载配置文件至目标参数
func LoadConfig(result interface{}) (err error) {
	bs, err := jsoniter.Marshal(defaultConfig.data)
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(bs, result)
	return
}

// 加载指定key到目标文件，key的形式为[a.b[0].c]
func LoadSpecifyConfig(key string, result interface{}) (err error) {
	value, err := defaultConfig.searchConfig(key)
	if err != nil {
		return
	}
	bs, err := jsoniter.Marshal(value)
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(bs, result)
	return
}

// 获取配置的环境变量
func GetRootPath() string {
	return defaultConfig.envConfig.RootPath
}

// 获取当前运行环境
func GetRunMode() RunModeType {
	return defaultConfig.envConfig.RunMode
}

// 获取指定配置项的interface{}类型值
func GetRawValue(key string) (ret interface{}, err error) {
	ret, err = defaultConfig.searchConfig(key)
	if err != nil {
		return
	}
	return
}

// 获取指定配置项的int类型值
func GetIntValue(key string) (ret int, err error) {
	value, err := defaultConfig.searchConfig(key)
	if err != nil {
		return
	}
	ret, isNumber := isNumber(value)
	if !isNumber {
		err = errors.New(fmt.Sprintf("[%v] can not convert to int type", value))
	}
	return
}

// 获取指定配置项的string类型值
func GetStringValue(key string) (ret string, err error) {
	value, err := defaultConfig.searchConfig(key)
	if err != nil {
		return
	}
	ret = fmt.Sprintf("%v", value)
	return
}

// 获取指定配置项的float64类型值
func GetFloat64Value(key string) (ret float64, err error) {
	value, err := defaultConfig.searchConfig(key)
	if err != nil {
		return
	}
	ret, isFloat64 := isFloat64(value)
	if !isFloat64 {
		err = errors.New(fmt.Sprintf("[%v] can not convert to float64 type", value))
	}
	return
}