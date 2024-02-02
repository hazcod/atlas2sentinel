package atlas

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func toJson(obj interface{}) (string, error) {
	switch obj.(type) {
	case string:
		return obj.(string), nil
	}

	b, err := json.Marshal(&obj)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ConvertLogtToMap(_ *logrus.Logger, logsBytes [][]byte) ([]map[string]string, error) {
	logs := make([]map[string]string, len(logsBytes))

	var err error

	for i, _ := range logs {
		cols := make(map[string]string)

		// TODO

		logs[i] = cols
	}

	return logs, err
}
