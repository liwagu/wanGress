package main

import (
	"os"
	"os/exec"
)

func main() {
	// 指定Envoy配置文件的路径
	configPath := "../../envoy/config/envoy.yaml"

	// 使用Docker启动Envoy进程
	cmd := exec.Command("docker", "run", "--rm", "-v", configPath+":/etc/envoy/envoy.yaml", "-p", "10000:10000", "envoyproxy/envoy:v1.29.2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
