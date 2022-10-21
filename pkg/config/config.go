package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	defaultAMIType      = "AL2_x86_64"
	defaultInstanceType = "m5.large"
	defaultRoleName     = "kubeprov-nodegroup-scope"
	defaultMaxSize      = int64(3)
	defaultDesiredState = int64(3)
	defaultMinSize      = int64(2)
)

type Config struct {
	Region       string `json:"region" yaml:"region" mapstructure:"region"`
	RoleName     string `json:"role_name" yaml:"role_name" mapstructure:"role_name"`
	ClusterName  string `json:"cluster_name" yaml:"cluster_name" mapstructure:"cluster_name"`
	AMIType      string `json:"ami_type" yaml:"ami_type" mapstructure:"ami_type"`
	InstanceType string `json:"instance_type" yaml:"instance_type" mapstructure:"instance_type"`
	DesiredState int64  `json:"desired_state" yaml:"desired_state" mapstructure:"desired_state"`
	MaxSize      int64  `json:"max_size" yaml:"max_size" mapstructure:"max_size"`
	MinSize      int64  `json:"min_size" yaml:"min_size" mapstructure:"min_size"`
}

func (c *Config) NodegroupName() string {
	return fmt.Sprintf("%s-nodegroup", c.ClusterName)
}

func InitAndParse() (*Config, error) {
	viper.AddConfigPath("$HOME/.kubeprov")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("ami_type", defaultAMIType)
	viper.SetDefault("instance_type", defaultInstanceType)
	viper.SetDefault("role_name", defaultRoleName)
	viper.SetDefault("max_size", defaultMaxSize)
	viper.SetDefault("min_size", defaultMinSize)
	viper.SetDefault("desired_state", defaultDesiredState)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if viper.Get("cluster_name") == "" {
				return nil, err
			}
		}
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
