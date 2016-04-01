package config

import (
	"regexp"
	"strconv"
	"strings"
)

type ConfigValueParser interface {
	ParseLine(string) interface{}
	ParseValue(string) interface{}
}

type ConfigIntParser struct{}
type ConfigStringParser struct{}
type ConfigUintParser struct{}
type ConfigFloatParser struct{}
type ConfigIntArrayParser struct{ sep string }
type ConfigStrArrayParser struct{ sep string }

func (parser *ConfigIntParser) ParseLine(value_str string) interface{} {
	r := regexp.MustCompile(`^[-]?[\d]+$`)
	suffix := r.FindString(value_str)
	if suffix != "" {
		return parser.ParseValue(suffix)
	}
	return nil
}

func (parser *ConfigIntParser) ParseValue(value_str string) interface{} {
	i, err := strconv.ParseInt(value_str, 0, 64)
	if err != nil {
		return 0
	}
	return i
}

func (parser *ConfigStringParser) ParseLine(value_str string) interface{} {
	return value_str
}

func (parser *ConfigStringParser) ParseValue(value_str string) interface{} {
	return value_str
}

func (parser *ConfigUintParser) ParseLine(value_str string) interface{} {
	r := regexp.MustCompile(`^[\d]+$`)
	suffix := r.FindString(value_str)
	if suffix != "" {
		return parser.ParseValue(suffix)
	}
	return nil
}

func (parser *ConfigUintParser) ParseValue(value_str string) interface{} {
	i, err := strconv.ParseUint(value_str, 0, 64)
	if err != nil {
		return nil
	}
	return i
}

func (parser *ConfigFloatParser) ParseLine(value_str string) interface{} {
	r := regexp.MustCompile(`^\d+\.?\d+|\d`)
	suffix := r.FindString(value_str)
	if suffix != "" {
		return parser.ParseValue(suffix)
	}
	return nil
}

func (parser *ConfigFloatParser) ParseValue(value_str string) interface{} {
	f, err := strconv.ParseFloat(value_str, 64)
	if err != nil {
		return nil
	}
	return f
}

func (parser *ConfigIntArrayParser) ParseLine(value_str string) interface{} {
	sp := strings.Split(value_str, parser.sep)
	var int_array []int64
	if len(sp) > 0 {
		for _, int_str := range sp {
			val := parser.ParseValue(int_str)
			if val != nil {
				int_array = append(int_array, val.(int64))
			}
		}
	}

	if len(int_array) > 0 {
		return int_array
	}
	return nil
}

func (parser *ConfigIntArrayParser) ParseValue(value_str string) interface{} {
	i, err := strconv.ParseInt(value_str, 10, 64)
	if err != nil {
		return nil
	}
	return i
}

func (parser *ConfigStrArrayParser) ParseLine(value_str string) interface{} {
	sp := strings.Split(value_str, parser.sep)
	if len(sp) > 0 {
		return sp
	}
	return nil
}

func (parser *ConfigStrArrayParser) ParseValue(value_str string) interface{} {
	//do nothing
	return nil
}
