package memory

type RecordRef struct {
	Kind string
	ID   string
}

func (r RecordRef) valid() bool {
	return r.Kind != "" && r.ID != ""
}
