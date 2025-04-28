package configs

import "github.com/spf13/viper"

type Conf struct {
	CepApiUrl     string `mapstructure:"cep_api_url"`
	WeatherApiUrl string `mapstructure:"weather_api_url"`
	WeatherApiKey string `mapstructure:"weather_api_key"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
