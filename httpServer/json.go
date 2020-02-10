package httpServer

import "encoding/json"

type arrStringMap map[string][]string

func (asm arrStringMap) MarshalJSON() ([]byte, error) {
	res := []byte{'{', '\n'}
	i := -1
	for key, arr := range asm {
		res = append(res, []byte{' ', ' ', ' ', ' '}...)
		res = append(res, '"')
		res = append(res, []byte(key)...)
		res = append(res, []byte("\": ")...)
		aByte, err := json.Marshal(arr)
		if err != nil {
			return nil, err
		}

		res = append(res, aByte...)
		i++
		if i != len(asm)-1 {
			res = append(res, ',')
		}
		res = append(res, '\n')
	}

	res = append(res, '}')
	return res, nil
}
