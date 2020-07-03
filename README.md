## config

用于读取配置文件，支持`yaml`和`json`格式的配置文件，主要的功能为:

1、根据环境变量`CONFIG_RUN_MODE`的设置读取不同的配置文件，可区分`local`、`dev`、`ci`、`test`、`uat`、`prod` 环境；

2、可以便捷读取配置文件值，例如通过`config.GetStringValue("redis.address[0]")`这样来读取配置文件；

3、可以通过设置环境变量的值来更改加载自配置文件的值，例如设置环境变量`config_redis_address="redis.address[0]=127.0.0.1:7000"`来更改配置，这样可以增加下环境变量，重启服务，而不用重新发版。

### 1、quick start

项目目录如下：

```
.
├── config
│   ├── ci.yml
│   ├── dev.yml
│   ├── local.yml
│   ├── prod.yml
│   ├── test.yml
│   └── uat.yml
└── main.go
```

1、查找指定环境的配置文件

查找配置文件的路径依次为：1）环境变量`CONFIG_ROOT_PATH`指定目录下的`config`目录；2）当前工作目录；3）当前工作目录的上级目录。通过`CONFIG_RUN_MODE`区分加载不同的配置文件，如`CONFIG_RUN_MODE="test"`会加载`test.yml`配置，支持的配置文件后缀有`yml`、`yaml`、`json`三种。

假定当前为`test`环境，config包加载`test.yml`配置文件：

```
# xxx/config/test.yml

service_name:
  test

redis:
  mode: cluster
  address:
    - 127.0.0.1:7000
    - 127.0.0.1:7001
    - 127.0.0.1:7002
    - 127.0.0.1:7003

num:
  1

float_num:
  1.11
```

2、读取配置文件

1) 将配置文件加载至结构化变量：

```go
type (
	configData struct {
		ServiceName string       `json:"service_name"`
		Redis       *RedisConfig `json:"redis"`
		Num         int          `json:"num"`
		FloatNum    float64      `json:"float_num"`
	}
	RedisConfig struct {
		Mode    string   `json:"mode"`
		Address []string `json:"address"`
	}
)

func main() {
	var c configData
	err := config.LoadConfig(&c)
	if err != nil {
		log.Fatal("load config error: ", err.Error())
	}
  
  var redisConfig RedisConfig
	err := config.LoadSpecifyConfig("redis", &redisConfig)
	if err != nil {
		log.Fatal("load redis config error: ", err.Error())
	}
}
```

2)读取指定的配置项值

```go
serviceName, _ := config.GetStringValue("service_name")
address, _ := config.GetRawValue("redis.address[0]")
```

3、通过更改环境变量来在项目启动时更改配置项

配置的环境变量须以`config`或者`CONFIG`开头，如：`config_service_name="service_name=test1"`，`config_redis_address="redis.address[0]=127.0.0.1:7000"`。

### 2、api list

- func LoadConfig(result interface{}) (err error)

- func LoadSpecifyConfig(key string, result interface{}) (err error)

- func GetRootPath() string
- func GetRunMode() RunModeType
- func GetRawValue(key string) (ret interface{}, err error)
- func GetIntValue(key string) (ret int, err error)
- func GetStringValue(key string) (ret string, err error)
- func GetFloat64Value(key string) (ret float64, err error)

