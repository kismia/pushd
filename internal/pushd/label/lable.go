package label

import (
	"errors"
)

var ErrInconsistentLabel = errors.New("inconsistent label")

type Label struct {
	Name  string
	Value string
}

func New(name, val string) Label {
	return Label{
		Name:  name,
		Value: val,
	}
}

type Labels []Label

func FromByteSlices(args [][]byte) (Labels, error) {
	if len(args)%2 != 0 {
		return nil, ErrInconsistentLabel
	}

	set := make(Labels, int(len(args)/2))

	li := 0

	for i, arg := range args {
		if i%2 == 0 {
			set[li] = Label{
				Name: string(arg),
			}
			li++
		} else {
			set[li-1].Value = string(arg)
		}
	}

	return set, nil
}

func (s Labels) Names() []string {
	names := make([]string, len(s))

	for i, l := range s {
		names[i] = l.Name
	}

	return names
}

func (s Labels) Values() []string {
	values := make([]string, len(s))

	for i, l := range s {
		values[i] = l.Value
	}

	return values
}
