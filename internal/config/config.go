package config

import (
	"github.com/spf13/viper"
)

//Config is viper config reader.
type Config struct {
	viper *viper.Viper
}

func New(name, ext, path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType(ext)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{viper: v}, nil
}

//Token gets bot token.
func (c *Config) Token() string {
	return c.viper.GetString("token")
}

//YtToken gets YouTube API token.
func (c *Config) YtToken() string {
	return c.viper.GetString("youtubeToken")
}

//AppID gets app ID.
func (c *Config) AppID() string {
	return c.viper.GetString("appID")
}

//FileDirectory gets directory for media files.
func (c *Config) FileDirectory() string {
	return c.viper.GetString("fileDir")
}

//LogLevel ge.
func (c *Config) LogLevel() string {
	return c.viper.GetString("logLevel")
}
