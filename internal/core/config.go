package core

import "github.com/null-bd/logger/types"

func defaultConfig() *types.Config {
	return &types.Config{
		ServiceName: "unknown",
		Environment: "development",
		LogLevel:    types.InfoLevel,
		Format:      "json",
		OutputPaths: []string{"stdout"},
	}
}
