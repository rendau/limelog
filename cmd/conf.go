package cmd

import (
	"strings"

	dopLogger "github.com/rendau/dop/adapters/logger"
	"github.com/rendau/dop/dopTools"
	"github.com/spf13/viper"
)

var conf = struct {
	Debug                  bool   `mapstructure:"DEBUG"`
	LogLevel               string `mapstructure:"LOG_LEVEL"`
	HttpListen             string `mapstructure:"HTTP_LISTEN"`
	HttpCors               bool   `mapstructure:"HTTP_CORS"`
	SwagHost               string `mapstructure:"SWAG_HOST"`
	SwagBasePath           string `mapstructure:"SWAG_BASE_PATH"`
	SwagSchema             string `mapstructure:"SWAG_SCHEMA"`
	MongoUsername          string `mapstructure:"MONGO_USERNAME"`
	MongoPassword          string `mapstructure:"MONGO_PASSWORD"`
	MongoHost              string `mapstructure:"MONGO_HOST"`
	MongoDbName            string `mapstructure:"MONGO_DB_NAME"`
	MongoReplicaSet        string `mapstructure:"MONGO_REPLICA_SET"`
	AuthPassword           string `mapstructure:"AUTH_PASSWORD"`
	SessionToken           string `mapstructure:"SESSION_TOKEN"`
	NfTelegramBotToken     string `mapstructure:"NF_TELEGRAM_BOT_TOKEN"`
	NfTelegramChatId       int64  `mapstructure:"NF_TELEGRAM_CHAT_ID"`
	NfTelegramLevels       string `mapstructure:"NF_TELEGRAM_LEVELS"`
	NfTelegramLevelsParsed []string
	InputGelfAddr          string `mapstructure:"INPUT_GELF_ADDR"`
	InputHttpAddr          string `mapstructure:"INPUT_HTTP_ADDR"`
}{}

func confLoad() {
	dopTools.SetViperDefaultsFromObj(conf)

	viper.SetDefault("DEBUG", "false")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("HTTP_LISTEN", ":80")
	viper.SetDefault("SWAG_HOST", "example.com")
	viper.SetDefault("SWAG_BASE_PATH", "/")
	viper.SetDefault("SWAG_SCHEMA", "https")
	viper.SetDefault("INPUT_GELF_ADDR", ":9234")
	viper.SetDefault("INPUT_HTTP_ADDR", ":9235")
	viper.SetDefault("MONGO_HOST", "localhost:27017")

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	_ = viper.Unmarshal(&conf)
}

func confParse(lg dopLogger.WarnAndError) {
	conf.NfTelegramLevelsParsed = confParseNfLevels(conf.NfTelegramLevels)
}

func confParseNfLevels(src string) []string {
	result := make([]string, 0)

	for _, lvl := range strings.Split(src, ",") {
		lvl = strings.ToLower(strings.TrimSpace(lvl))
		if lvl != "" {
			result = append(result, lvl)
		}
	}

	return result
}
