package std

import "strings"

type Bridge []string

func (source Bridge) String() string {
	return strings.Join(source, " ⇝ ")
}
