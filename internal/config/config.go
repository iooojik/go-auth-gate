package config

import (
	"fmt"
	"os"

	"github.com/iooojik/go-auth-gate/internal/repository/session"
	"github.com/iooojik/go-auth-gate/pkg/apple"
	"github.com/iooojik/go-auth-gate/pkg/google"
	"github.com/iooojik/go-auth-gate/pkg/jwt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AppleSignIn  apple.Config   `yaml:"appleSignIn"`
	GoogleSignIn google.Config  `yaml:"googleSignIn"`
	JWT          jwt.Config     `yaml:"jwt"`
	SQL          session.Config `yaml:"sql"`
}

func Load(path string) Config {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("open config file: %s", err))
	}

	dec := yaml.NewDecoder(f)

	dec.KnownFields(true)

	cfg := new(Config)

	err = dec.Decode(cfg)
	if err != nil {
		panic(fmt.Sprintf("decode config file: %s", err))
	}

	return *cfg
}
