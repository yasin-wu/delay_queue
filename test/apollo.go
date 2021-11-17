package test

import "github.com/apolloconfig/agollo/v4/env/config"

var apolloConf = &config.AppConfig{
	AppID:          "SampleApp",
	Cluster:        "dev",
	IP:             "http://localhost:8080",
	NamespaceName:  "application",
	IsBackupConfig: true,
	Secret:         "49792384a8e14a999e13df1f7aa064fe",
}
