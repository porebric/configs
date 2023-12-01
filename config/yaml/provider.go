package yaml

import (
	"context"
	"fmt"
	"io"

	"github.com/porebric/configs/config"
	"github.com/porebric/configs/keys"
	"github.com/porebric/logger"

	yaml "gopkg.in/yaml.v3"
)

type provider struct {
	values map[string]*config.Value
}

type valuesDataType struct {
	Value map[string]string `yaml:"values"`
}

func Init(ctx context.Context, r io.Reader, k map[string]keys.ConfigType) (*provider, error) {
	yamlData := new(valuesDataType)
	if err := yaml.NewDecoder(r).Decode(yamlData); err != nil {
		return nil, fmt.Errorf("yaml: cannot decode values yaml: %w", err)
	}

	p := &provider{
		values: make(map[string]*config.Value),
	}

	for name, key := range k {
		strValue, ok := yamlData.Value[name]
		if !ok {
			logger.Debug(ctx, "yaml: value does not exist in yaml", "key", name)
			if val := config.Convert(ctx, key.Default, key.Type); val != nil {
				p.values[name] = config.New(val, name)
			}
			continue
		}
		if val := config.Convert(ctx, strValue, key.Type); val != nil {
			p.values[name] = config.New(val, name)
		}
	}

	return p, nil
}

func (y *provider) Value(key string) *config.Value {
	val, ok := y.values[key]
	if ok {
		return val
	}
	return nil
}

func (y *provider) Watch(_ context.Context, key string, _ ...config.WatchFn) *config.Value {
	return y.Value(key)
}
