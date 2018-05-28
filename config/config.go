package config

import (
	"strings"
	"fmt"
	"github.com/spf13/viper"
)

const fileName = "config"

var (
	initFlg bool = false
)

func InitConfig() error {
	viper.SetEnvPrefix(fileName)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	//viper.SetConfigName(cmdRoot)
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error when reading %s config file:%s", fileName, err)
	}

	fmt.Println("配置文件加载成功:", fileName)
	initFlg = true
	return nil
}

func HasModuleInit() bool {
	return initFlg
}
