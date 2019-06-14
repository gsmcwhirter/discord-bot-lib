package bot

// JSONMarshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}
