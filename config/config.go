package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var err error

type SituationType string
type ENV string //deactivated

const (
	LoadingConfigFile SituationType = "LoadingConfigFile"
	Unmarshalling     SituationType = "Unmashalling"
	// Env               ENV           = "DEV" // deactivated
)

type ConfigError struct {
	Situation SituationType
	Message   error
}

type Config struct {
	Shcool_id string `mapstructure:"LOILO_DISTRICT_ID"`
	User_id   string `mapstructure:"LOILO_DISTRICT_USER_ID"`
	User_pw   string `mapstructure:"LOILO_DISTRICT_USER_PW"`
	// ENV       ENV    `mapstructure:"ENV"` // deactivated
}

func (c ConfigError) Error() string {
	return fmt.Sprintf("Error situation: %s - message: %s", c.Situation, c.Message)
}

func init() {
	viper.AddConfigPath("config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
}

/* Deactivated
func checkOnProd(c *Config) bool {
	// 設定ファイルのENVフィールドで判定
	// 認証情報が空かつENV=PROD以外はダメ（PRODではユーザーに入力させる想定）
	return c.Shcool_id != "" && c.User_id != "" && c.User_pw != "" && c.ENV != "PROD"
}
*/

// envファイルがあったらそれを使う
func Load() (*Config, error) {
	err = viper.ReadInConfig()
	if err != nil {
		return nil, ConfigError{Situation: LoadingConfigFile, Message: fmt.Errorf("error occured in reading config file: %w", err)}
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, ConfigError{Situation: Unmarshalling, Message: fmt.Errorf("error occured in unmarshalling: %w", err)}
	}
	return &cfg, nil
}
