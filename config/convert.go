package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/porebric/logger"
)

func Convert(ctx context.Context, v string, t string) any {

	switch t {
	case TypeString:
		return v

	case TypeInt:
		value, err := strconv.Atoi(v)
		if err == nil {
			return value
		}
		logger.Error(ctx, fmt.Errorf("strconv.Atoi: %w", err), "convert error")

	case TypeFloat:
		value, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return value
		}
		logger.Error(ctx, fmt.Errorf("strconv.ParseFloat: %w", err), "convert error")

	case TypeBool:
		value, err := strconv.ParseBool(v)
		if err == nil {
			return value
		}

		logger.Error(ctx, fmt.Errorf("strconv.ParseBool: %w", err), "convert error")

	case TypeDuration:
		value, err := time.ParseDuration(v)
		if err == nil {
			return value
		}

		logger.Error(ctx, fmt.Errorf("time.ParseDuration: %w", err), "convert error")

	case TypeIntMap:
		value := make(map[string]int)
		elems := strings.Split(v, ",")
		var err error
		for _, elem := range elems {
			e := strings.Split(elem, ":")
			if len(e) != 2 {
				err = errors.New("invalid map data: " + v)
				break
			}
			value[e[0]], err = strconv.Atoi(e[1])
			if err != nil {
				break
			}
		}
		if err == nil {
			return value
		}
		logger.Error(ctx, fmt.Errorf("parse map: %w", err), "convert error")
	case TypeStringSlice:
		elems := strings.Split(v, ",")
		value := make([]string, len(elems))
		var err error
		for i, elem := range elems {
			value[i] = elem
		}
		if err == nil {
			if len(value) == 0 {
				logger.Warn(ctx, "empty string slice", "value", v)
			}
			return value
		}
		logger.Error(ctx, fmt.Errorf("parse map: %w", err), "convert error")
	default:
		logger.Warn(ctx, "unknown value item type "+t)
	}

	return nil
}
