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
