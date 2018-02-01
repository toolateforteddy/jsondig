package jsondig

// Package jsondig lets you easily dig into a decoded json object
// to grab values deep within, without having to make all of the
// intermediate structs.

import (
	"fmt"
	"strings"
)

func JsonDig(v interface{}, path ...string) (interface{}, error) {
	if v == nil {
		return nil, &DigError{
			path: path[0:0],
		}
	}
	if len(path) == 0 {
		return v, nil
	}
	if msi, ok := v.(map[string]interface{}); ok {
		val, err := JsonDig(msi[path[0]], path[1:]...)
		if err != nil {
			if e, ok := err.(*DigError); ok {
				if e.v == nil {
					e.v = v
				}
				e.path = append(e.path, path[0])
			}
			return nil, err
		}
		return val, nil
	}
	return nil, &DigError{
		path: path[0:0],
		v:    v,
	}
}

type DigError struct {
	path []string
	v    interface{}
}

func (d *DigError) Error() string {
	numPaths := len(d.path)
	reverse := make([]string, numPaths)
	for i, p := range d.path {
		reverse[numPaths-1-i] = p
	}
	pathStr := strings.Join(reverse, ".")
	successFullPathStr := strings.Join(reverse[0:numPaths-1], ".")

	return fmt.Sprintf(
		"Could not find object at %q.\nFound %#v at %q.",
		pathStr,
		d.v,
		successFullPathStr,
	)
}
