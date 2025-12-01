package hashmap

import "encoding/json"

// MarshalJSON implements json.Marshaler
func (m HashMap[K, V]) MarshalJSON() ([]byte, error) {
	result := make(map[K]V, len(m.data))
	for k, e := range m.data {
		if e.IsOccupied() {
			result[k] = *e.Value()
		}
	}
	return json.Marshal(result)
}

// UnmarshalJSON implements json.Unmarshaler
func (m *HashMap[K, V]) UnmarshalJSON(data []byte) error {
	var tmp map[K]V
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*m = New[K, V]()
	for k, v := range tmp {
		m.Insert(k, v)
	}
	return nil
}
