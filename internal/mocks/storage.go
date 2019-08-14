package mocks

type MockStorage struct {
	data map[string]map[string][]byte
}

func NewMockStorage() MockStorage {
	return MockStorage{data: map[string]map[string][]byte{}}
}

func (_ MockStorage) Close() {

}

func (m MockStorage) GetAll(collection string) (map[string]string, error) {
	value, exists := m.data[collection]
	if !exists {
		return map[string]string{}, nil
	}

	result := map[string]string{}
	for k, v := range value {
		result[k] = string(v)
	}

	return result, nil
}

func (m MockStorage) GetAllRaw(collection string) (map[string][]byte, error) {
	value, exists := m.data[collection]
	if !exists {
		return map[string][]byte{}, nil
	}

	return value, nil
}

func (m MockStorage) Put(collection string, key string, value []byte) error {
	_, e := m.data[collection]
	if e == false {
		m.data[collection] = map[string][]byte{}
	}
	m.data[collection][key] = value

	return nil
}

func (m MockStorage) CreateCollectionIfNotExist(collection string) {

}
