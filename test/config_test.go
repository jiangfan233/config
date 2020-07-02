package test

import (
	"github.com/jiangfan233/config"
	"testing"
)

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

func TestLoadConfig(t *testing.T) {
	var c configData
	err := config.LoadConfig(&c)
	if err != nil {
		t.Fatal("load config error: ", err.Error())
	}
	if c.ServiceName != "test" {
		t.Fatal("load service config failed")
	}
	if c.Redis.Address[0] != "127.0.0.1:7000" {
		t.Fatal("load redis address config failed")
	}
	if c.Num != 1 {
		t.Fatal("load num config failed")
	}
	if c.FloatNum != 1.11 {
		t.Fatal("load float_num config failed")
	}
}

func TestLoadSpecifyConfig(t *testing.T) {
	var c RedisConfig
	err := config.LoadSpecifyConfig("redis", &c)
	if err != nil {
		t.Fatal("load redis config error: ", err.Error())
	}
	if len(c.Address) > 0 && (c.Address[0] != "127.0.0.1:7000" || c.Mode != "cluster") {
		t.Fatal("redis config is not the same with config file")
	}
}

func TestGetRunMode(t *testing.T) {
	runMode := config.GetRunMode()
	if runMode != "local" {
		t.Fatal("get run mode error")
	}
}

func TestGetRawValue(t *testing.T) {
	serviceName, err := config.GetRawValue("service_name")
	if err != nil {
		t.Fatal("get service_name config error: ", err.Error())
	}
	if serviceName != "test" {
		t.Fatalf("service_name is %s now, but actual is \"test\"", serviceName)
	}
	address, err := config.GetRawValue("redis.address[0]")
	if err != nil {
		t.Fatal("get redis.address[0] config error: ", err.Error())
	}
	if address != "127.0.0.1:7000" {
		t.Fatalf("redis.address[0] is %s now, but actual is \"127.0.0.1:7000\"", address)
	}
}

func TestGetIntValue(t *testing.T) {
	num, err := config.GetIntValue("num")
	if err != nil {
		t.Fatal("get num config error: ", err.Error())
	}
	if num != 1 {
		t.Fatalf("num is %d now, but actual is 1", num)
	}
}

func TestGetStringValue(t *testing.T) {
	s, err := config.GetStringValue("service_name")
	if err != nil {
		t.Fatal("get service_name config error: ", err.Error())
	}
	if s != "test" {
		t.Fatalf("service_name is %s now, but actual is \"test\"", s)
	}
}

func TestGetFloat64Value(t *testing.T) {
	f, err := config.GetFloat64Value("float_num")
	if err != nil {
		t.Fatal("get float_num config error: ", err.Error())
	}
	if f != float64(1.11) {
		t.Fatalf("float_num is %v now, but actual is 1.11", f)
	}
}

func TestReplaceConfig(t *testing.T) {
	err := config.ReplaceConfig("service_name", "test1")
	if err != nil {
		t.Fatal("replace service_name config error: ", err.Error())
	}
	s, _ := config.GetStringValue("service_name")
	if s != "test1" {
		t.Fatal("replace service_name config failed")
	}
	err = config.ReplaceConfig("redis.address[0]", "127.0.0.1:7004")
	if err != nil {
		t.Fatal("replace redis.address[0] config error: ", err.Error())
	}
	s, _ = config.GetStringValue("redis.address[0]")
	if s != "127.0.0.1:7004" {
		t.Fatal("replace redis.address[0] config failed")
	}
}
