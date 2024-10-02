package dvrip

import (
	"net"
	"time"
)

var magicEnd = [2]byte{0x0A, 0x00}

type Client struct {
	conn net.Conn

	Session  uint32
	Sequence uint32

	lastPing time.Duration
}

type Header struct {
	Head           uint8  // B
	Version        uint8  // B
	_              byte   // padding
	_              byte   // padding
	Session        uint32 // I
	SequenceNumber uint32 // I
	_              byte   // padding
	_              byte   // padding
	MsgID          uint16 // H
	Len            uint32 // I
}

type Response struct {
	Name      string `json:"Name"`
	Ret       int    `json:"Ret"`
	SessionID string `json:"SessionID"`
}

const (
	CodeLogin     uint16 = 1000
	CodeKeepAlive uint16 = 1006

	CodeSystemInfo uint16 = 1020
	CodeConfigSet  uint16 = 1040

	CodeSystemManager             = 1450
	CodeAuthorityList      uint16 = 1470
	CodeUsers              uint16 = 1472
	CodeSystemUpgrade      uint16 = 1525
	CodeStartSystemUpgrade uint16 = 0x5F0
	CodeSendFile           uint16 = 0x5F2
)

const (
	statusOK                                  int = 100
	statusUnknownError                        int = 101
	statusUnsupportedVersion                  int = 102
	statusRequestNotPermitted                 int = 103
	statusUserAlreadyLoggedIn                 int = 104
	statusUserIsNotLoggedIn                   int = 105
	statusUsernameOrPasswordIsIncorrect       int = 106
	statusUserDoesNotHaveNecessaryPermissions int = 107
	statusPasswordIsIncorrect                 int = 203
	statusStartOfUpgrade                      int = 511
	statusUpgradeWasNotStarted                int = 512
	statusUpgradeDataErrors                   int = 513
	statusUpgradeError                        int = 514
	statusUpgradeSuccessful                   int = 515
)

var mappedStatusCodes = map[int]string{
	statusOK:                                  "OK",
	statusUnknownError:                        "Unknown error",
	statusUnsupportedVersion:                  "Unsupported version",
	statusRequestNotPermitted:                 "Request not permitted",
	statusUserAlreadyLoggedIn:                 "User already logged in",
	statusUserIsNotLoggedIn:                   "User is not logged in",
	statusUsernameOrPasswordIsIncorrect:       "Username or password is incorrect",
	statusUserDoesNotHaveNecessaryPermissions: "User does not have necessary permissions",
	statusPasswordIsIncorrect:                 "Password is incorrect",
	statusStartOfUpgrade:                      "Start of upgrade",
	statusUpgradeWasNotStarted:                "Upgrade was not started",
	statusUpgradeDataErrors:                   "Upgrade data errors",
	statusUpgradeError:                        "Upgrade error",
	statusUpgradeSuccessful:                   "Upgrade successful",
}

var mappedOpcodes = map[uint16]string{
	CodeSystemUpgrade:      "OPSystemUpgrade",
	CodeUsers:              "Users",
	CodeAuthorityList:      "AuthorityList",
	CodeSystemInfo:         "SystemInfo",
	CodeStartSystemUpgrade: "OPSystemUpgrade",
}
