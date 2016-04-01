package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

const (
	CFG_INT = iota
	CFG_STRING
	CFG_UINT64
	CFG_FLOAT64
	CFG_INT_ARRAY
	CFG_STRING_ARRAY
)

const DEFAULT_ARRAY_SEP = "|"

type Config struct {
	once      sync.Once
	values    map[string]interface{}
	array_sep string
}

func (c *Config) importLoader(v ...interface{}) []*ConfigLoader {
	var loaders []*ConfigLoader
	idx := 0
	sep := DEFAULT_ARRAY_SEP

	if c.array_sep != "" {
		sep = c.array_sep
	}

	for idx < len(v) {
		//todo: check prefix

		loader := new(ConfigLoader)
		loader.prefix = v[idx].(string)
		idx++
		loader.config_type = v[idx].(int)
		idx++
		loader.defalut_value = v[idx].(string)
		idx++

		switch loader.config_type {
		case CFG_INT:
			loader.value_parser = &ConfigIntParser{}
		case CFG_STRING:
			loader.value_parser = &ConfigStringParser{}
		case CFG_UINT64:
			loader.value_parser = &ConfigUintParser{}
		case CFG_FLOAT64:
			loader.value_parser = &ConfigFloatParser{}
		case CFG_INT_ARRAY:
			loader.value_parser = &ConfigIntArrayParser{sep: sep}
		case CFG_STRING_ARRAY:
			loader.value_parser = &ConfigStrArrayParser{sep: sep}
		}
		loaders = append(loaders, loader)
	}
	return loaders
}

func (c *Config) LoadFromFile(filepath string, v ...interface{}) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		err = c.LoadFromBytes(nil, v...)
		if err != nil {
			return err
		} else {
			return errors.New("open file failed, use default config value")
		}
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		err = c.LoadFromBytes(nil, v...)
		if err != nil {
			return err
		} else {
			return errors.New("read file failed, use default config value")
		}
	}

	return c.LoadFromBytes(content, v...)
}

func (c *Config) LoadFromBytes(content []byte, v ...interface{}) (err error) {
	c.once.Do(func() {
		c.values = make(map[string]interface{})
	})

	if len(v)%3 != 0 {
		return errors.New("Load config arguments num error!")
	}

	loaders := c.importLoader(v...)

	if len(content) > 0 {
		sp := bytes.Split(content, []byte("\n"))
		for _, line := range sp {
			for _, loader := range loaders {
				loader.TryLoadValue(string(line))
			}
		}
	}

	for _, loader := range loaders {
		if loader.matched {
			c.values[loader.prefix] = loader.value
		} else {
			c.values[loader.prefix] = loader.value_parser.ParseValue(loader.defalut_value)
		}
	}

	return nil
}

func (c *Config) SetArraySep(sep string) {
	c.array_sep = sep
}

func (c *Config) GetValue(key string) interface{} {
	value, _ := c.values[key]
	return value
}

func (c *Config) GetString(key string) string {
	value, _ := c.values[key]
	return value.(string)
}

func (c *Config) GetInt64(key string) int64 {
	value, _ := c.values[key]
	return value.(int64)
}

func (c *Config) GetUint64(key string) uint64 {
	value, _ := c.values[key]
	return value.(uint64)
}

func (c *Config) GetFloat64(key string) float64 {
	value, _ := c.values[key]
	return value.(float64)
}

func (c *Config) GetIntArray(key string) []int64 {
	value, _ := c.values[key]
	return value.([]int64)
}

func (c *Config) GetStringArray(key string) []string {
	value, _ := c.values[key]
	return value.([]string)
}
