package script

type Errors struct {
	List []error
}

func (t *Errors) Add(err error) {
	if err != nil {
		t.List = append(t.List, err)
	}
}

func (t *Errors) First() error {
	if len(t.List) > 0 {
		return t.List[0]
	}
	return nil
}

func (t *Errors) IsEmpty(err error) bool {
	return len(t.List) == 0
}

func (t *Errors) Clear() {
	t.List = nil
}
