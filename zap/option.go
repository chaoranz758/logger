// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option interface {
	apply(cfg *config)
}

type ExtraKey string

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type CoreConfig struct {
	Enc zapcore.Encoder
	Ws  zapcore.WriteSyncer
	Lvl zapcore.LevelEnabler
}

type config struct {
	extraKeys   []ExtraKey
	coreConfigs []CoreConfig
	zapOpts     []zap.Option
}

// defaultCoreConfig default zapcore config: json encoder, atomic level, stdout write syncer
func defaultCoreConfig() *CoreConfig {
	// default log encoder
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// default log level
	lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	// default write syncer stdout
	ws := zapcore.AddSync(os.Stdout)

	return &CoreConfig{
		Enc: enc,
		Ws:  ws,
		Lvl: lvl,
	}
}

// defaultConfig default config
func defaultConfig() *config {
	return &config{
		coreConfigs: []CoreConfig{*defaultCoreConfig()},
		zapOpts:     []zap.Option{},
	}
}

// WithCoreEnc zapcore encoder
func WithCoreEnc(enc zapcore.Encoder) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs[0].Enc = enc
	})
}

// WithCoreWs zapcore write syncer
func WithCoreWs(ws zapcore.WriteSyncer) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs[0].Ws = ws
	})
}

// WithCoreLevel zapcore log level
func WithCoreLevel(lvl zap.AtomicLevel) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs[0].Lvl = lvl
	})
}

// WithCores zapcore
func WithCores(coreConfigs ...CoreConfig) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs = coreConfigs
	})
}

// WithZapOptions add origin zap option
func WithZapOptions(opts ...zap.Option) Option {
	return option(func(cfg *config) {
		cfg.zapOpts = append(cfg.zapOpts, opts...)
	})
}

// WithExtraKeys allow you log extra values from context
func WithExtraKeys(keys []ExtraKey) Option {
	return option(func(cfg *config) {
		for _, k := range keys {
			if !inArray(k, cfg.extraKeys) {
				cfg.extraKeys = append(cfg.extraKeys, k)
			}
		}
	})
}
