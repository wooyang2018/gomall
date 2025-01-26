// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"net"
)

// GetLocalIPv4 获取本地IPv4地址
// 返回值：
//   - string: 本地IPv4地址
//   - error: 获取过程中发生的错误
func GetLocalIPv4() (string, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		// 如果获取失败，抛出错误
		panic(err)
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 检查网络接口是否为非回环接口且已启用
		if iface.Flags&net.FlagLoopback != net.FlagLoopback && iface.Flags&net.FlagUp != 0 {
			// 获取网络接口的所有地址
			addrs, err := iface.Addrs()
			if err != nil {
				// 如果获取失败，继续下一个网络接口
				continue
			}

			// 遍历所有地址
			for _, addr := range addrs {
				// 将地址转换为IPNet类型
				if ipNet, ok := addr.(*net.IPNet); ok {
					// 检查IP地址是否为非回环地址且为IPv4地址
					if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
						// 返回IPv4地址
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}

	// 如果没有找到IPv4地址，返回错误
	return "", errors.New("get local IP error")
}

func MustGetLocalIPv4() string {
	ipv4, err := GetLocalIPv4()
	if err != nil {
		panic("get local IP error")
	}
	return ipv4
}
