package normalize

import (
	"encoding/json"
	"errors"
)

func DecodeBatch(body []byte) ([]json.RawMessage, error) {
	var single json.RawMessage
	if json.Unmarshal(body, &single) != nil {
		return nil, errors.New("Expected JSON object, array, or AWS Records envelope")
	}
	var batch []json.RawMessage
	if len(single) > 0 && single[0] == '[' {
		if err := json.Unmarshal(single, &batch); err != nil {
			return nil, err
		}
	} else {
		var envelope map[string]json.RawMessage
		if json.Unmarshal(single, &envelope) != nil || envelope == nil {
			return nil, errors.New("Expected JSON object or array")
		}
		if records, ok := envelope["Records"]; ok {
			if err := json.Unmarshal(records, &batch); err != nil {
				return nil, err
			}
			for i, item := range batch {
				var m map[string]any
				if json.Unmarshal(item, &m) != nil || m == nil {
					return nil, errors.New("AWS Records must contain objects")
				}
				m["source"] = "aws"
				batch[i], _ = json.Marshal(m)
			}
		} else {
			batch = []json.RawMessage{single}
		}
	}
	if len(batch) < 1 || len(batch) > 1000 {
		return nil, errors.New("Batch must contain 1–1000 logs")
	}
	return batch, nil
}
