package pint

import "strings"

type tagOptions map[string]string

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], genTagOptions(tag[idx+1:])
	}
	return tag, tagOptions{}
}

func genTagOptions(options string) tagOptions {
	o := tagOptions{}
	parts := strings.Split(options, ",")
	for _, p := range parts {
		if idx := strings.Index(p, ":"); idx != -1 {
			o[p[:idx]] = p[idx+1:]
		} else {
			o[p] = ""
		}
	}
	return o
}

func (o tagOptions) contains(option string) bool {
	_, ok := o[option]
	return ok
}

func (o tagOptions) get(option string) (string, bool) {
	val, ok := o[option]
	return val, ok
}
