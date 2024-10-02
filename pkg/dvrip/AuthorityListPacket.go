package dvrip

import (
	"errors"
	"fmt"
)

type test struct{}

func (c *Client) GetAuthorityList() (*test, error) {
	var result = new(test)

	_, resp, err := c.Instruct(CodeAuthorityList, "", nil)
	if err != nil {
		return nil, errors.New("failed to send instruction to device: " + err.Error())
	}

	fmt.Println(string(resp))

	return result, nil
}
