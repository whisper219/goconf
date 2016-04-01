package config

import (
	"regexp"
)

type ConfigLoader struct {
	matched       bool
	prefix        string
	config_type   int
	defalut_value string
	value         interface{}
	value_parser  ConfigValueParser
}

func (loader *ConfigLoader) TryLoadValue(line string) {
	if loader.matched {
		return
	}

	var cur_str string = line
	r := regexp.MustCompile(`^[\s]*` + loader.prefix + `[\s]+`)
	prefix_idx := r.FindStringIndex(cur_str)
	if prefix_idx == nil {
		return
	}

	cur_str = cur_str[prefix_idx[1]:]
	ret := loader.value_parser.ParseLine(cur_str)
	if ret != nil {
		loader.value = ret
		loader.matched = true
	}
}
