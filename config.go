package main

var (
	userName = "blockcloud"
	password = "blockcloud2018"
)

type Config struct {
	jenkinsConfig JenkinsConfig
}

type JenkinsConfig struct{
	userName string
	password string
}

func NewConfig() Config{
	config := Config{}

	jenkinsconfig := JenkinsConfig{}
	jenkinsconfig.userName = userName
	jenkinsconfig.password = password

	config.jenkinsConfig = jenkinsconfig

	return config
}