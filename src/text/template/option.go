// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the code to handle template options.

package template

import "strings"

// missingKeyAction defines how to respond to indexing a map with a key that is not present.
type missingKeyAction int

const (
	mapInvalid   missingKeyAction = iota // Return an invalid reflect.Value.
	mapZeroValue                         // Return the zero value for the map element.
	mapError                             // Error out
)

type option struct {
	missingKey missingKeyAction
}

// onPanicAction defines how to handle template function panics.
type onPanicAction int

const (
	panicRecover   onPanicAction = iota // Recover the panic and return it as an error
	panicNop                            // Do nothing, let the caller handle the panic
)

type option struct {
	missingKey missingKeyAction
	onPanic    onPanicAction
}


// Option sets options for the template. Options are described by
// strings, either a simple string or "key=value". There can be at
// most one equals sign in an option string. If the option string
// is unrecognized or otherwise invalid, Option panics.
//
// Known options:
//
// missingkey: Control the behavior during execution if a map is
// indexed with a key that is not present in the map.
//
//	"missingkey=default" or "missingkey=invalid"
//		The default behavior: Do nothing and continue execution.
//		If printed, the result of the index operation is the string
//		"<no value>".
//	"missingkey=zero"
//		The operation returns the zero value for the map type's element.
//	"missingkey=error"
//		Execution stops immediately with an error.
//
// onpanic: Control the behavior during execution if a panic occurs.
//
//	"onpanic=recover"
//		The default behavior: Recover the panic and return it as an error.
//	"onpanic=nop"
//		Do nothing, let the caller handle the panic.
func (t *Template) Option(opt ...string) *Template {
	t.init()
	for _, s := range opt {
		t.setOption(s)
	}
	return t
}

func (t *Template) setOption(opt string) {
	if opt == "" {
		panic("empty option string")
	}
	// key=value
	if key, value, ok := strings.Cut(opt, "="); ok {
		switch key {
		case "missingkey":
			switch value {
			case "invalid", "default":
				t.option.missingKey = mapInvalid
				return
			case "zero":
				t.option.missingKey = mapZeroValue
				return
			case "error":
				t.option.missingKey = mapError
				return
			}
		case "onpanic":
			switch value {
			case "recover":
				t.option.onPanic = panicRecover
				return
			case "nop":
				t.option.onPanic = panicNop
				return
			}
		}
	}
	panic("unrecognized option: " + opt)
}
