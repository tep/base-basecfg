// Copyright 2019 Timothy E. Peoples
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

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
