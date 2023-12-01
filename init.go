package configs

import (
	"context"
	"github.com/porebric/configs/config"
	"github.com/porebric/configs/config/env"
	"io"

	"github.com/porebric/configs/config/yaml"
	"github.com/porebric/configs/keys"
	"github.com/porebric/logger"
)

var keysStorage keyTypes
var providerStorage *Storage

type Storage struct {
	keysReader        io.Reader
	yamlConfigsReader io.Reader

	configsProviders map[keys.SourceType]config.Provider
}

type keyTypes struct {
	yamlConfigs map[string]bool
	envConfigs  map[string]bool
}

func New() *Storage {
	return new(Storage)
}

func (s *Storage) YamlConfigs(r io.Reader) *Storage {
	s.yamlConfigsReader = r
	return s
}

func (s *Storage) KeysReader(r io.Reader) *Storage {
	s.keysReader = r
	return s
}

func (s *Storage) Init(ctx context.Context) error {
	configKeys, err := keys.Init(s.keysReader)
	if err != nil {
		return err
	}
	s.configsProviders = make(map[keys.SourceType]config.Provider)

	keysStorage.yamlConfigs = make(map[string]bool)
	if s.yamlConfigsReader != nil {
		// init yaml config keys
		for k, _ := range configKeys[keys.Yaml].Configs {
			keysStorage.yamlConfigs[k] = true
		}
		s.configsProviders[keys.Yaml], err = yaml.Init(ctx, s.yamlConfigsReader, configKeys[keys.Yaml].Configs)
	}

	keysStorage.envConfigs = make(map[string]bool)
	for k, _ := range configKeys[keys.Env].Configs {
		keysStorage.envConfigs[k] = true
	}
	s.configsProviders[keys.Env], err = env.Init(ctx, configKeys[keys.Env].Configs)

	providerStorage = s
	return err
}

func Value(ctx context.Context, key string) *config.Value {
	if _, ok := keysStorage.yamlConfigs[key]; ok {
		return providerStorage.configsProviders[keys.Yaml].Value(key)
	}
	if _, ok := keysStorage.envConfigs[key]; ok {
		return providerStorage.configsProviders[keys.Env].Value(key)
	}

	logger.Warn(ctx, "unknown value key", "key", key)

	return nil
}

func Watch(ctx context.Context, key string, fns ...config.WatchFn) *config.Value {
	if _, ok := keysStorage.yamlConfigs[key]; ok {
		return providerStorage.configsProviders[keys.Yaml].Watch(ctx, key, fns...)
	}
	logger.Warn(ctx, "unknown value key", "key", key)

	return nil
}
