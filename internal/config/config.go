package config

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func GetConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if err = cfg.AuthConf.Validate(); err != nil {
		return nil, fmt.Errorf("authConf validate error: %w", err)
	}
	return &cfg, nil
}

type Config struct {
	App        *AppCfg     `yaml:"app"`
	Repository *Repository `yaml:"repository"`
	AuthConf   *AuthConf   `yaml:"auth"`
	Facility   string      `yaml:"facility"`
}

type AuthConf struct {
	PrivateKey           string        `yaml:"private_key"`
	PublicKey            string        `yaml:"public_key"`
	PrivateKeyByte       []byte        `yaml:"-"`
	PublicKeyByte        []byte        `yaml:"-"`
	AccessKeyExpiration  time.Duration `yaml:"access_key_expiration"`
	RefreshKeyExpiration time.Duration `yaml:"refresh_key_expiration"`
}

func (a *AuthConf) Validate() error {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(a.PublicKey)
	if err != nil {
		return err
	}
	a.PublicKeyByte = decodedPublicKey

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(a.PrivateKey)
	if err != nil {
		return err
	}
	a.PrivateKeyByte = decodedPrivateKey

	return nil
}

type MongoRepo struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
}

type AppCfg struct {
	Port    int           `yaml:"port"`
	RTO     time.Duration `yaml:"rto"`
	WTO     time.Duration `yaml:"wto"`
	AdmConf *AdminConf    `yaml:"admin"`
}

type AdminConf struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}
type Repository struct {
	Mongo *MongoRepo `yaml:"mongo"`
}
