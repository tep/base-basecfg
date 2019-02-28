package basecfg

import (
	"reflect"
	"strings"
)

func (fd *featureDefn) extractDefaults() {
	f := fd.Feature
	v := reflect.Indirect(reflect.ValueOf(f))
	t := reflect.TypeOf(f).Elem()

	if v.Kind() != reflect.Struct || t.Kind() != reflect.Struct {
		return
	}

	fd.defaults = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		if fi := getFieldInfo(t, i); fi != nil {
			fd.defaults[fd.label.Key(fi.key)] = v.FieldByName(fi.name).Interface()
		}
	}
}

type fieldInfo struct {
	name string
	key  string
}

func getFieldInfo(t reflect.Type, i int) *fieldInfo {
	sf := t.Field(i)

	// skip if unexported (PkgPath is populated only for unexported fields)
	if sf.PkgPath != "" {
		return nil
	}

	tag := sf.Tag.Get("cfg")
	if tag == "" {
		return nil
	}

	parts := strings.Split(tag, ",")
	key := parts[0]

	if len(parts) > 1 {
		for _, p := range parts[1:] {
			if p == "nodefault" {
				return nil
			}
		}
	}

	return &fieldInfo{sf.Name, key}
}
