package dvrip

import (
	"errors"
)

func (c *Client) Reboot() error {
	_, _, err := c.Instruct(CodeSystemManager, "OPMachine", map[string]interface{}{
		"Action": "Reboot",
	})

	if err != nil {
		return errors.New("failed to send instruction to device: " + err.Error())
	}

	return nil
}
