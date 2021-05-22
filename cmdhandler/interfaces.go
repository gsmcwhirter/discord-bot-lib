package cmdhandler

// JSONMarshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type JSONMarshaler interface {
	MarshalToJSON() ([]byte, error) // yes, this is intentionally different than stdlib
}
