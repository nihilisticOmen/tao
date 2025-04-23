package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

// Server 结构体: 存储服务名称、地址、版本和权重信息
type Server struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`    //服务地址
	Version string `json:"version"` //服务版本
	Weight  int64  `json:"weight"`  //服务权重
}

// BuildPrefix : 构建服务在etcd中的前缀路径/<服务名>/[<版本>/]
func BuildPrefix(info Server) string {
	if info.Version == "" {
		return fmt.Sprintf("/%s/", info.Name)
	}
	return fmt.Sprintf("/%s/%s/", info.Name, info.Version)
}

// BuildRegPath : 构建服务完整注册路径/<服务名>/[<版本>/]<地址>
func BuildRegPath(info Server) string {
	return fmt.Sprintf("%s%s", BuildPrefix(info), info.Addr)
}

// ParseValue : 将JSON数据反序列化为Server结构
func ParseValue(value []byte) (Server, error) {
	info := Server{}
	if err := json.Unmarshal(value, &info); err != nil {
		return info, err
	}
	return info, nil
}

// SplitPath : 从etcd路径中提取服务地址信息
func SplitPath(path string) (Server, error) {
	info := Server{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return info, errors.New("invalid path")
	}
	info.Addr = strs[len(strs)-1]
	return info, nil
}

// Exist /Remove : 地址列表操作辅助函数
func Exist(l []resolver.Address, addr resolver.Address) bool {
	for i := range l {
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

// Remove helper function
func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func BuildResolverUrl(app string) string {
	return schema + ":///" + app
}
