package basecfg

import "errors"

type FeatureError error

var (
	ErrRegistrationClosed = FeatureError(errors.New("feature registration is closed"))
	ErrMissingOneOfName   = FeatureError(errors.New("cannot register OneOf without a name"))
)

type DuplicateLabelError struct {
	Dupe Label
	error
}

type MultipleOneOfError struct {
	Name     string
	Feature1 Label
	Feature2 Label
	error
}

func multipleOneOfError(name string, feat1, feat2 Label) *MultipleOneOfError {
	if feat2 < feat1 {
		feat1, feat2 = feat2, feat1
	}

	return &MultipleOneOfError{
		Name:     name,
		Feature1: feat1,
		Feature2: feat2,
		error:    errors.New("multiple configurations found for mutually exclusive feature set"),
	}
}
