package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name  string `mapstructure:"name"`
	Mode  string `mapstructure:"mode"`
	Verso string `mapstructure:"version"`
	Port  int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Password  string `mapstructure:"password"`
	Port      int    `mapstructure:"port"`
	Db        int    `mapstructure:"db"`
	PoolsSize int    `mapstructure:"pools_size"`
}

func Init() (err error) {
	viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	//viper.SetConfigType("yaml")   // 指定配置文件类型(专用于从远程获取配置信息时使用)
	viper.AddConfigPath(".")   // 指定查找配置文件的路径
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到Conf变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return
	}

	viper.WatchConfig() // 监控配置文件变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
			return
		}
	})
	return

}
