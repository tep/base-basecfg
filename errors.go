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
