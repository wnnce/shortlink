package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

var (
	globalConfig *Bootstrap
	readerSet    map[*ConfigureReader]struct{}
)

type Bootstrap struct {
	Server Server `json:"server" yaml:"server"`
	Etcd   Etcd   `json:"etcd" yaml:"etcd"`
}

type Server struct {
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	ServiceName string `json:"serviceName" yaml:"service-name"`
}

type Etcd struct {
	Endpoints   []string      `json:"endpoints" yaml:"endpoints"`
	DialTimeout time.Duration `json:"dialTimeout" yaml:"dial-timeout"`
	Username    string        `json:"username" yaml:"username"`
	Password    string        `json:"password" yaml:"password"`
}

// ConfigureReader 配置读取函数
// 会下发所有的配置的复制体，避免配置被修改
// 返回的func()用于在服务被关闭时执行一些指定的清理操作
// error则是根据配置进行初始化时的错误情况
type ConfigureReader func(bootstrap Bootstrap) (func(), error)
type ReaderList []ConfigureReader

func NewReaderList(readers ...ConfigureReader) ReaderList {
	return readers
}

// InitConfig 初始化配置
// path 传入配置文件路径
func InitConfig(path string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(dir, path)
	file, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	bootstrap := &Bootstrap{}
	if err = yaml.Unmarshal(file, bootstrap); err != nil {
		return err
	}
	globalConfig = bootstrap
	return nil
}

// RegisterConfigureReader 注册配置文件读取函数
func RegisterConfigureReader(readerLists ...ReaderList) {
	if readerSet == nil {
		readerSet = make(map[*ConfigureReader]struct{})
	}
	for _, readerList := range readerLists {
		for _, reader := range readerList {
			if _, ok := readerSet[&reader]; !ok {
				readerSet[&reader] = struct{}{}
			}
		}
	}
}

// IssueConfigure 分发配置，供各个函数读取配置并执行初始化操作
func IssueConfigure() func() {
	if readerSet == nil || globalConfig == nil {
		return nil
	}
	cleans := make([]func(), 0)
	for reader, _ := range readerSet {
		clean, err := (*reader)(*globalConfig)
		if err != nil {
			panic(err)
		}
		if clean != nil {
			cleans = append(cleans, clean)
		}
	}
	// 返回资源清理闭包
	return func() {
		for _, clean := range cleans {
			clean()
		}
	}
}

func DefaultConfig() Bootstrap {
	return *globalConfig
}
