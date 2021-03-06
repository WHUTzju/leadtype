// Copyright 2016 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package ltml

type HasParent interface {
	Parent() interface{}
	SetParent(value interface{}) error
}
