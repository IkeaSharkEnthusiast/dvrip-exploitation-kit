package dvrip

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Client) ConfigSet(name string, data map[string]interface{}) error {
	_, resp, err := c.Instruct(CodeConfigSet, name, map[string]interface{}{
		name: data,
	})

	if err != nil {
		return errors.New("failed to send instruction to device: " + err.Error())
	}

	if err := json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data["Ret"].(float64) != float64(statusOK) {
		return fmt.Errorf("unexpected status code: %s", mappedStatusCodes[int(data["Ret"].(float64))])
	}

	return nil
}
