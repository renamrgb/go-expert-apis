// Package main
package main

import "github.com/renamrgb/go-expert-apis/configs"

func main() {
	config, _ := configs.LoadConfig(".")

	println(config.DBConfig.DBDriver)
}
