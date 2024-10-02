package dvrip

import (
	"encoding/json"
	"errors"
	"fmt"
)

type UsersResponse struct {
	*Response
	Users []user `json:"Users"`
}

type user struct {
	AuthorityList []string    `json:"AuthorityList"`
	Group         string      `json:"Group"`
	Memo          string      `json:"Memo"`
	Name          string      `json:"Name"`
	NoMD5         interface{} `json:"NoMD5"`
	Password      string      `json:"Password"`
	Reserved      bool        `json:"Reserved"`
	Sharable      bool        `json:"Sharable"`
}

func (c *Client) GetUsers() (*UsersResponse, error) {
	var result = new(UsersResponse)

	_, resp, err := c.Instruct(CodeUsers, "", nil)
	if err != nil {
		return nil, errors.New("failed to send instruction to device: " + err.Error())
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	if result.Ret != statusOK {
		return nil, fmt.Errorf("unexpected status code; %s", mappedStatusCodes[result.Ret])
	}

	return result, nil
}
