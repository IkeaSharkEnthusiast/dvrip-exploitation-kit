package dvrip

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type LoginResponse struct {
	AliveInterval int    `json:"AliveInterval"`
	ChannelNum    int    `json:"ChannelNum"`
	DeviceType    string `json:"DeviceType "`
	ExtraChannel  int    `json:"ExtraChannel"`
	Ret           int    `json:"Ret"`
	SessionID     string `json:"SessionID"`
}

func (c *Client) Login(username, password string) error {
	var res = new(LoginResponse)

	body, err := json.Marshal(map[string]string{
		"EncryptType": "MD5",
		"LoginType":   "DVRIP-WEB",
		"UserName":    username,
		"PassWord":    password,
	})

	if err != nil {
		return err
	}

	err = c.Write(CodeLogin, body, uint32(len(body))+2)
	if err != nil {
		return err
	}

	_, resp, err := c.Read(true)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resp, &res); err != nil {
		return err
	}

	if res.Ret != statusOK {
		return fmt.Errorf("unexpected status code: %s", mappedStatusCodes[res.Ret])
	}

	c.lastPing = time.Duration(res.AliveInterval) * time.Second
	session, err := strconv.ParseUint(res.SessionID, 0, 32)
	if err != nil {
		return err
	}

	c.Session = uint32(session)
	return nil
}
