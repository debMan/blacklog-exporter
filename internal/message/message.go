package message

type Config struct {
	DataJSONPath string `json:"data_json_path,omitempty" koanf:"level"`
}

// type Message map[string]interface{} // TODO

type Message []byte
