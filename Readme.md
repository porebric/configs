# Library for working with configs

The library is designed for parsing a file with configs.

The library supports the following types of configs:
- string
- int
- float
- bool
- duration
- int_map
- string_slice

## How to use it

### Creating two yaml files.

- In the first one, we denote all keys types and default values for <b>configs_keys.yml</b>

```yml
config:
  db_host:
    type: string
    default: "localhost"
    source: yaml
  db_name:
    type: string
    default: "test"
    source: yaml
  db_user:
    type: string
    default: "root"
    source: yaml
  db_password:
    type: string
    default: "qwerty"
    source: yaml
  db_port:
    type: int
    default: "5432"
    source: yaml
```

- In the second one, the real values of these configs are <b>configs.yml</b>

```yml
values:
  db_host: localhost
  db_name: prod
  db_user: user
  db_password: 12345678
  db_port: 5432
```
### Example

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/porebric/configs"
)

const (
	configsKeysPath = "./configs_keys.yml"
	configsPath     = "./configs.yml"
)

func main() {
	ctx := context.Background()
	keysReader, err := os.Open(configsKeysPath)
	if err != nil {
		panic(err)
	}

	confReader, err := os.Open(configsPath)
	if err != nil {
		panic(err)
	}

	if err = configs.New().KeysReader(keysReader).YamlConfigs(confReader).Init(ctx); err != nil {
		panic(err)
	}

	fmt.Println(configs.Value(ctx, "db_host").String())
	fmt.Println(configs.Value(ctx, "db_user").String())
        fmt.Println(configs.Value(ctx, "db_password").String())
        fmt.Println(configs.Value(ctx, "db_name").String())
        fmt.Println(configs.Value(ctx, "db_port").Int())
}
```

<b>P.S.</b> The supported types and values can be seen in config_test.go