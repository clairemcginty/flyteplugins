// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package config

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (Config) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (Config) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in Config and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg Config) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("Config", pflag.ExitOnError)
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "environment"), defaultConfig.Environment.String(), "Environment endpoint for Presto to use")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "defaultRoutingGroup"), defaultConfig.DefaultRoutingGroup, "Default Presto routing group")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "defaultUser"), defaultConfig.DefaultUser, "Default Presto user")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "rateLimiter.name"), defaultConfig.RateLimiter.Name, "The name of the rate limiter")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "rateLimiter.syncPeriod"), defaultConfig.RateLimiter.SyncPeriod.String(), "The duration to wait before the cache is refreshed again")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "rateLimiter.workers"), defaultConfig.RateLimiter.Workers, "Number of parallel workers to refresh the cache")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "rateLimiter.lruCacheSize"), defaultConfig.RateLimiter.LruCacheSize, "Size of the cache")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "rateLimiter.metricScope"), defaultConfig.RateLimiter.MetricScope, "The prefix in Prometheus used to track metrics related to Presto")
	return cmdFlags
}
