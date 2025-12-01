package slice

import "encoding/json"

// MarshalJSON implements json.Marshaler
func (s Slice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.values)
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	var values []T
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	*s = Init(values...)
	return nil
}
