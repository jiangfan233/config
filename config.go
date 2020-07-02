package config

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type (
	envConfiguration struct {
		RunMode  RunModeType `default:"local" split_words:"true"`
		RootPath string      `split_words:"true"`
	}

	configuration struct {
		envConfig    *envConfiguration
		data      map[interface{}]interface{}
	}
)

const configDir = "/config/"

type RunModeType string

// 运行环境
const (
	RunModeLocal RunModeType = "local"
	RunModeDev   RunModeType = "dev"
	RunModeCi    RunModeType = "ci"
	RunModeTest  RunModeType = "test"
	RunModeUat   RunModeType = "uat"
	RunModeProd  RunModeType = "prod"
)

// 配置文件类型
const (
	jsonType = "json"
	yamlType = "yaml"
	ymlType  = "yml"
)

var defaultConfig configuration

var configNotExist = errors.New("config item not exist")

func init() {
	var envConfig envConfiguration
	// 检查环境变量中的当前运行环境，并加载相关环境变量配置
	if err := envconfig.Process("config", &envConfig); err != nil {
		panic("no CONFIG_RUN_MODE exported")
	}
	defaultConfig.envConfig = &envConfig
	// 加载项目目录下的 json/yaml 格式的配置文件
	defaultConfig.loadConfigFile()
	// 替换环境变量中以CONFIG_开头的配置项
	defaultConfig.setConfigPrefixEnvData()
}

// 加载配置文件
func (c *configuration) loadConfigFile() {
	if c == nil {
		c = new(configuration)
	}
	if c.envConfig.RunMode == "" {
		panic("project run mode is not specified")
	}
	configFilePath, fileType := findConfigFileByPriority(c.envConfig.RootPath, c.envConfig.RunMode)
	fileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic("load config file error")
	}
	c.data = make(map[interface{}]interface{})
	if fileType == jsonType {
		// 清除注释
		reg, _ := regexp.Compile("/\\*.+?\\*/")
		fileString := string(fileBytes)
		fileString = reg.ReplaceAllString(fileString, "")
		if err := jsoniter.Unmarshal([]byte(fileString), &c.data); err != nil {
			panic(errors.New(fmt.Sprintf("config file: %s is not json format", configFilePath)))
		}
	}
	if err = yaml.Unmarshal([]byte(fileBytes), &c.data); err != nil {
		panic(errors.New(fmt.Sprintf("config file: %s is not yaml format", configFilePath)))
	}
}

// 替换环境变量中以CONFIG_开头的配置项
func (c *configuration) setConfigPrefixEnvData() {
	for _, item := range os.Environ() {
		l := strings.Split(item, "=")
		if len(l) == 3 && (strings.HasPrefix(l[0], "CONFIG") || (strings.HasPrefix(l[0], "config"))) {
			err := c.replaceConfig(l[1], l[2])
			if err != nil {
				log.Printf("replace env config prefix key <%s> failed", l[0])
			}
		}
	}
}

// 按优先级查找配置文件
// 查找顺序 环境变量CONFIG_ROOT_PATH的config文件夹 > 当前目录的config文件夹 > 上级目录的config文件夹
func findConfigFileByPriority(rootPath string, runMode RunModeType) (configFilePath string, fileType string) {
	if runMode == "" {
		return
	}
	currentDir, _ := os.Getwd()
	searchDirs := []string{rootPath, currentDir, currentDir + "/.."}
	if rootPath == "" {
		searchDirs = []string{currentDir, currentDir + "/.."}
	}
	for _, dir := range searchDirs {
		if configFilePath, fileType = searchConfigFile(dir+configDir, runMode); configFilePath != "" && fileType != "" {
			return
		}
	}
	if configFilePath == "" || fileType == "" {
		panic("no config file found")
	}
	log.Println("find by priority, file path is: ", configFilePath)
	return
}

func searchConfigFile(rootPath string, runMode RunModeType) (filePath string, fileType string) {
	list := []string{jsonType, yamlType, ymlType}
	for _, fileType = range list {
		filePath = rootPath + string(runMode) + "." + fileType
		if fileExists(filePath) {
			return
		}
	}
	return "", ""
}

// 替换掉configuration对象中的key的值
// eg：
// 	configuration.rawData={"foo": 0}} key = "foo", value=1 -> {"foo": 1}
// 	configuration.rawData={"foo": {"bar": 0}} key = "foo.bar", value=1 -> {"foo": {"bar": 1}}
// 	configuration.rawData={"foo": {"bar": [0, 1, 2]}} key = "foo.bar[0]", value=1 -> {"foo": {"bar": [1, 1, 2]}}
// 	configuration.rawData={"foo": {"bar": [0, 1, 2]}} key = "foo.bar[0]", value=1 -> {"foo": {"bar": [1, 1, 2]}}
// 替换数组字段时，支持数组内对象字段替换，嵌套对象字段为数组也可以替换，如 a.b[o].c[1]
func (c *configuration) replaceConfig(key string, value interface{}) (err error) {
	list := splitKey(key)
	length := len(list)
	var currentData interface{} = c.data
	var previousKeys []string
	for i, k := range list {
		previousKeys = append(previousKeys, k)
		currentKey := composeKey(previousKeys) // 记录查找到的位置，方便错误时打印输出
		index, isNumber := isNumber(k)
		if !isNumber {
			var valueExists bool
			dataInterfaceMap, isInterfaceMap := currentData.(map[interface{}]interface{})
			if !isInterfaceMap {
				dataStringMap, isStringMap := currentData.(map[string]interface{})
				if !isStringMap {
					log.Printf("search config [%s] unreachable", currentKey)
					return configNotExist
				}
				currentData, valueExists = dataStringMap[k]
				if !valueExists {
					log.Printf("search config [%s] not exist", currentKey)
					return configNotExist
				}
				if i+1 == length {
					dataStringMap[k] = value
					log.Printf("replace config [%s] success, the current value is [%v]", key, value)
					return
				}
			} else {
				currentData, valueExists = dataInterfaceMap[k]
				if !valueExists {
					log.Printf("search config [%s] not exist", currentKey)
					return configNotExist
				}
				if i+1 == length {
					dataInterfaceMap[k] = value
					log.Printf("replace config [%s] success, the current value is [%v]", key, value)
					return
				}
			}
		} else {
			data, ok := currentData.([]interface{})
			if !ok {
				log.Printf("replace config %s[%d] unreachable", currentKey, index)
				return configNotExist
			}
			if length <= index {
				log.Printf("replace config %s[%d] not exist", currentKey, index)
				return configNotExist
			}
			if i+1 == length {
				data[index] = value
				log.Printf("replace config [%s] success, the current value is [%v]", key, value)
				return
			}
			currentData = data[index]
		}
	}
	return
}

// 与replaceConfig方法逻辑基本一致，只是将替换改成返回值
func (c *configuration) searchConfig(key string) (ret interface{}, err error) {
	list := splitKey(key)
	length := len(list)
	ret = c.data
	var previousKeys []string
	for i, k := range list {
		previousKeys = append(previousKeys, k)
		currentKey := composeKey(previousKeys)
		index, isNumber := isNumber(k)
		if !isNumber {
			var valueExists bool
			dataInterfaceMap, isInterfaceMap := ret.(map[interface{}]interface{})
			if !isInterfaceMap {
				dataStringMap, isStringMap := ret.(map[string]interface{})
				if !isStringMap {
					log.Printf("search config [%s] unreachable", currentKey)
					return nil, configNotExist
				} else {
					ret, valueExists = dataStringMap[k]
				}
			} else {
				ret, valueExists = dataInterfaceMap[k]
			}
			if !valueExists {
				log.Printf("search config [%s] not exist", currentKey)
				return nil, configNotExist
			}
			if i+1 == length {
				log.Printf("search config [%s] success, the value is [%v]", key, ret)
				return
			}
		} else {
			data, ok := ret.([]interface{})
			if !ok {
				log.Printf("search config %s[%d] unreachable", currentKey, index)
				return nil, configNotExist
			}
			if length <= index {
				log.Printf("search config %s[%d] not exist", currentKey, index)
				return nil, configNotExist
			}
			ret = data[index]
			if i+1 == length {
				log.Printf("search config [%s] success, the value is [%v]", key, ret)
				return
			}
		}
	}
	return
}

// 提供单测用
func ReplaceConfig(key string, value interface{}) (err error) {
	if defaultConfig.envConfig.RunMode != RunModeLocal {
		err = errors.New("this function is only used in local mode for test")
		return
	}
	err = defaultConfig.replaceConfig(key, value)
	return
}