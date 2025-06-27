package config

import (
	"github.com/spf13/viper"
)

type ProjectConfig struct {
	Name string `mapstructure:"name"`
}

type ApiConfig struct {
	Port       string `mapstructure:"port"`
	MaxNum     int    `mapstructure:"max_num"`
	SessionTTL int    `mapstructure:"session_ttl"`
}

type LogConfig struct {
	Compress    bool   `mapstructure:"compress"`
	LeepDays    int    `mapstructure:"leep_days"`
	Level       string `mapstructure:"level"`
	Mode        string `mapstructure:"mode"`
	Path        string `mapstructure:"path"`
	ServiceName string `mapstructure:"service_name"`
}

type KvConf struct {
	Redis []*RedisConfig `toml:"redis" mapstructure:"redis" json:"redis"`
}

type RedisConfig struct {
	MasterName string `toml:"master_name" mapstructure:"master_name" json:"master_name"`
	Pass       string `mapstructure:"pass"`
	Host       string `mapstructure:"host"`
	Type       string `mapstructure:"type"`
}

type DBConfig struct {
	Database           string `mapstructure:"database"`
	Host               string `mapstructure:"host"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	Port               int    `mapstructure:"port"`
	MaxOpenConns       int    `mapstructure:"max_open_conns"`
	LogLevel           string `mapstructure:"log_level"`
	MaxConnMaxLifetime int    `mapstructure:"max_conn_max_lifetime"`
	MaxIdleConns       int    `mapstructure:"max_idle_conns"`
}

type ChainSupported struct {
	Name     string `mapstructure:"name"`
	ChainID  int    `mapstructure:"chain_id"`
	Endpoint string `mapstructure:"endpoint"`
}

type EasySwapMarket struct {
	ApiKey   string `mapstructure:"apikey"`
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	Contract string `mapstructure:"contract"`
	Fee      int    `mapstructure:"fee"`
}

type ImageConfig struct {
	ValidFileTypes     []string `mapstructure:"valid_file_type"`
	Timeout            int      `mapstructure:"time_out"`
	PublicIPFSGateways []string `mapstructure:"public_ipfs_gateways"`
	LocalIPFSGateways  []string `mapstructure:"local_ipfs_gateways"`
	DefaultOSSUri      string   `mapstructure:"default_oss_uri"`
}

type MetadataParse struct {
	NameTags       []string `mapstructure:"name_tags"`
	ImageTags      []string `mapstructure:"image_tags"`
	AttributesTags []string `mapstructure:"attributes_tags"`
	TraitNameTags  []string `mapstructure:"trait_name_tags"`
	TraitValueTags []string `mapstructure:"trait_value_tags"`
}

type Config struct {
	Project  ProjectConfig    `mapstructure:"project_cfg"`
	API      ApiConfig        `mapstructure:"api"`
	Log      LogConfig        `mapstructure:"log"`
	Kv       *KvConf          `toml:"kv" json:"kv"`
	DB       DBConfig         `mapstructure:"db"`
	Chains   []ChainSupported `mapstructure:"chain_supported"`
	EasySwap EasySwapMarket   `mapstructure:"easyswap_market"`
	Image    ImageConfig      `mapstructure:"image_cfg"`
	Metadata MetadataParse    `mapstructure:"metadata_parse"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return defaultConfig(), nil
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return defaultConfig(), nil
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return defaultConfig(), nil
	}

	return &config, nil
}

func defaultConfig() *Config {
	return &Config{
		Project: ProjectConfig{
			Name: "MetaFarm",
		},
		API: ApiConfig{
	Port:       ":80",
	MaxNum:     500,
	SessionTTL: 86400,
},
		Log: LogConfig{
			Compress:    false,
			LeepDays:    7,
			Level:       "info",
			Mode:        "console",
			Path:        "logs/v1-backend",
			ServiceName: "v1-backend",
		},
		Kv: &KvConf{
			Redis: []*RedisConfig{{
				Pass: "123456",
				Host: "127.0.0.1:6379",
				Type: "node",
			}}},
		DB: DBConfig{
			Database:           "meta_farm",
			Host:               "127.0.0.1",
			User:               "meta_farm",
			Password:           "1qaz!QAZ",
			Port:               3306,
			MaxOpenConns:       1500,
			LogLevel:           "info",
			MaxConnMaxLifetime: 300,
			MaxIdleConns:       10,
		},
		Chains: []ChainSupported{{
			Name:     "sepolia",
			ChainID:  11155111,
			Endpoint: "https://rpc.ankr.com/eth_sepolia",
		}},
		EasySwap: EasySwapMarket{
			ApiKey:   "",
			Name:     "EasySwap",
			Version:  "1",
			Contract: "0x1466ceE9XXXXXXXXXXXXXXXXXXXcD4",
			Fee:      100,
		},
		Image: ImageConfig{
			ValidFileTypes:     []string{".jpeg", ".gif", ".png", ".mp4", ".jpg", ".glb", ".gltf", ".mp3", ".wav", ".svg"},
			Timeout:            40,
			PublicIPFSGateways: []string{"https://gateway.pinata.cloud/ipfs/", "https://cf-ipfs.com/ipfs/", "https://ipfs.infura.io/ipfs/", "https://ipfs.pixura.io/ipfs/", "https://ipfs.io/ipfs/", "https://www.via0.com/ipfs/"},
			LocalIPFSGateways:  []string{"https://gateway.pinata.cloud/ipfs/", "https://cf-ipfs.com/ipfs/", "https://ipfs.infura.io/ipfs/", "https://ipfs.pixura.io/ipfs/", "https://ipfs.io/ipfs/", "https://www.via0.com/ipfs/"},
			DefaultOSSUri:      "https://test.easyswap.link/",
		},
		Metadata: MetadataParse{
			NameTags:       []string{"name", "title"},
			ImageTags:      []string{"image", "image_url", "animation_url", "media_url", "image_data", "imageUrl"},
			AttributesTags: []string{"attributes", "properties", "attribute"},
			TraitNameTags:  []string{"trait_type"},
			TraitValueTags: []string{"value"},
		},
	}
}
