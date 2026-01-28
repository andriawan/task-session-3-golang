package config

type Template struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	DB struct {
		Host             string `mapstructure:"host"`
		Port             int    `mapstructure:"port"`
		User             string `mapstructure:"user"`
		Password         string `mapstructure:"password"`
		ConnectionString string `mapstructure:"connection_string"`
	} `mapstructure:"db"`
}
