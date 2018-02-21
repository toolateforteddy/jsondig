package jsondig

// Package jsondig lets you easily dig into a decoded json object
// to grab values deep within, without having to make all of the
// intermediate structs.

import (
	"fmt"
	"strconv"
	"strings"
)

func cleanArrayInd(ind string) (int, error) {
	return strconv.Atoi(strings.Trim(ind, "[]"))
}

func JsonDig(v interface{}, path ...string) (interface{}, error) {
	retVal := v
	lastVal := v
	pathInd := 0
	for {
		switch vv := retVal.(type) {
		case string, bool, float64, nil:
			retVal = vv
			if len(path) == pathInd && retVal != nil {
				return retVal, nil
			}
			return nil, &digError{
				path: path[0:pathInd],
				v:    lastVal,
			}
		case map[string]interface{}:
			if len(path) == pathInd && retVal != nil {
				return retVal, nil
			}

			lastVal = retVal
			retVal = vv[path[pathInd]]
			pathInd++
		case []interface{}:
			if len(path) == pathInd && retVal != nil {
				return retVal, nil
			}

			index, err := cleanArrayInd(path[pathInd])
			if err != nil {
				return nil, err
			}
			if len(vv) <= index {
				return nil, &digError{
					path: path[0 : pathInd+1],
					v:    vv,
				}
			}
			lastVal = retVal
			retVal = vv[index]
			pathInd++
		default:
			return nil, &digError{
				path: path[0:pathInd],
				v:    lastVal,
			}
		}
		if len(path) < pathInd {
			return nil, &digError{
				path: path[0 : pathInd-1],
				v:    lastVal,
			}
		}
	}
}

type digError struct {
	path []string
	v    interface{}
}

func (d *digError) Error() string {
	pathStr := strings.Join(d.path, ".")
	successFullPathStr := strings.Join(d.path[:len(d.path)-1], ".")

	return fmt.Sprintf(
		"Could not find object at %q.\nFound %#v at %q.",
		pathStr,
		d.v,
		successFullPathStr,
	)
}
