package env

import (
	"context"
	"os"
	"strings"

	"github.com/porebric/configs/config"
	"github.com/porebric/configs/keys"
	"github.com/porebric/logger"
)

type provider struct {
	values map[string]*config.Value
}

func Init(ctx context.Context, k map[string]keys.ConfigType) (*provider, error) {
	p := &provider{
		values: make(map[string]*config.Value),
	}

	for name, key := range k {
		strValue, ok := os.LookupEnv(strings.ToUpper(name))
		if !ok {
			logger.Debug(ctx, "env: value does not exist in env", "key", name)
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
