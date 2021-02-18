package script

type Error struct {
	First error
}

func (t *Error) Add(err error) {
	if t.First == nil && err != nil {
		t.First = err
	}
}

func (t *Error) IsNil() bool {
	return t.First == nil
}
