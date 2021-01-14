package atomic

// String is an atomic type-safe wrapper for string values.
type String struct {
	v Value
}

var _zeroString string

// NewString creates a new String.
func NewString(v string) *String {
	x := &String{}
	if v != _zeroString {
		x.Store(v)
	}
	return x
}

// Load atomically loads the wrapped string.
func (s *String) Load() string {
	if v := s.v.Load(); v != nil {
		return v.(string)
	}
	return _zeroString
}

// Store atomically stores the passed string.
func (s *String) Store(v string) {
	s.v.Store(v)
}

// String returns the wrapped value.
func (s *String) String() string {
	return s.Load()
}
