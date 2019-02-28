package basecfg

import (
	"github.com/spf13/pflag"
)

// Feature represents a configuration feature. The base type for a Feature should
// be a struct whose fields have a `cfg` tag.
type Feature interface {
	FlagSet(*pflag.FlagSet)
	// Defaults() map[string]interface{}
	Validate() error
}

// FeatureFunc is a function that returns a newly created Feature and should be
// the second argument to `Register()`.
type FeatureFunc func() Feature

// Register will register a new Feature with the global, in-memory Feature
// registry -- thus making it part of the current application's configuration
// set.
//
// It is often conventional for features to register themselves (in their own
// `init()` function) so that they're enabled implicitly by importing the
// Feature package.
//
// A nil error is returned on successful registration. If called after Features
// have been reified, ErrRegistrationClosed is returned. If the provided label
// has already been registered, an error of type DuplicateLabelError is
// returned.
func Register(l Label, f FeatureFunc) error {
	return registry.add(&featureDefn{label: l, create: f})
}

//TODO(tep): Allow caller to specify whether a OneOf is mandatory or optional

// RegisterOneOf is similar to Register in that it may be used to add a Feature
// to the global registry, but also takes a `oneOf` parameter to mark this
// Feature as a member of a mutual exclusion set. For each Feature registered
// with the same `oneOf` value, only one may be configured at runtime. If more
// than one of these Features is configured, validation will fail.
//
// If the value fo `oneOf` is the empty string, RegisterOneOf will return
// ErrMissingOneOfName -- otherwise, returned errors are as described for
// Register.
func RegisterOneOf(oneOf string, l Label, f FeatureFunc) error {
	if oneOf == "" {
		return ErrMissingOneOfName
	}

	return registry.add(&featureDefn{label: l, oneof: oneOf, create: f})
}