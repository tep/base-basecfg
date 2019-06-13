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
	"errors"
	"sort"
	"sync"
)

var registry featureRegistry

type featureRegistry struct {
	defs    map[Label]*featureDefn
	reified bool
	sync.Mutex
}

func (r *featureRegistry) add(fd *featureDefn) error {
	r.Lock()
	defer r.Unlock()

	if r.reified {
		return ErrRegistrationClosed
	}

	if r.defs == nil {
		r.defs = make(map[Label]*featureDefn)
	}

	if _, x := r.defs[fd.label]; x {
		return FeatureError(&DuplicateLabelError{fd.label, errors.New("duplicate feature label")})
	}

	r.defs[fd.label] = fd

	return nil
}

func (r *featureRegistry) reify() []*featureDefn {
	r.Lock()
	defer r.Unlock()

	r.reified = true

	if r.defs == nil || len(r.defs) == 0 {
		return nil
	}

	list := make([]*featureDefn, len(r.defs))

	for i, lbl := range r.labels() {
		f := r.defs[lbl]
		f.reify()
		list[i] = f
	}

	return list
}

func (r *featureRegistry) labels() []Label {
	list := make([]Label, len(r.defs))
	var i int
	for lbl := range r.defs {
		list[i] = lbl
		i++
	}

	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })

	return list
}

// A featureDefn is the result of registering a Feature. Prior to reification,
// defaults and Feature are nil
type featureDefn struct {
	label    Label
	oneof    string
	create   FeatureFunc
	defaults map[string]interface{}
	Feature
}

func (fd *featureDefn) reify() {
	fd.Feature = fd.create()
	fd.extractDefaults()
}
