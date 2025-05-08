package encoding

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	//UseNumber:              true,
}.Froze()

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalToString(v any) (string, error) {
	return json.MarshalToString(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func UnmarshalFromString(str string, v any) error {
	return json.UnmarshalFromString(str, v)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
