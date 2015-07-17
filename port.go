package generictypes

import (
	"github.com/juju/errgo"

	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func MustParseDockerPort(port string) DockerPort {
	var result DockerPort
	if err := parseDockerPort(port, &result); err != nil {
		panic(err.Error())
	}
	return result
}

func ParseDockerPort(port string) (DockerPort, error) {
	var result DockerPort
	if err := parseDockerPort(port, &result); err != nil {
		return result, errgo.Mask(err)
	}
	return result, nil
}

type modePortJSONFormat int

const (
	modePortJsonDocker modePortJSONFormat = 0
	modePortJsonNumber modePortJSONFormat = 1
	modePortJsonString modePortJSONFormat = 2
)

const (
	ProtocolTCP = "tcp"
	ProtocolUDP = "udp"
)

type DockerPort struct {
	// The port number.
	Port string

	// The protocol to use. "tcp" or "udp"
	Protocol string

	// How to format this port when marshalling as JSON.
	// 0 = format as string - ("port/protocol")
	// 1 = format as int - port
	// 2 = format as short port - "<port>"
	//
	// This is needed, because we need to marshal our ports the way we parsed them.
	// Otherwise the diff check in CheckForUnknownFields() would trigger when we
	// marshal `6379` as `"6379/tcp"`.
	formatJsonMode modePortJSONFormat
}

func (d DockerPort) String() string {
	return fmt.Sprintf("%s/%s", d.Port, d.Protocol)
}

func (d DockerPort) MarshalJSON() ([]byte, error) {
	switch d.formatJsonMode {
	case modePortJsonDocker:
		return json.Marshal(d.String())
	case modePortJsonNumber:
		if d.Protocol != ProtocolTCP {
			return nil, errgo.Newf("Invalid protocol for formatJsonMode=number")
		}
		i, err := strconv.Atoi(d.Port)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		return json.Marshal(i)
	case modePortJsonString:
		if d.Protocol != ProtocolTCP {
			return nil, errgo.Newf("Invalid protocol for formatJsonMode=number")
		}
		return json.Marshal(d.Port)
	default:
		panic("Invalid 'formatJsonMode'")
	}
}

func (d *DockerPort) UnmarshalJSON(data []byte) error {
	wasNumber := false
	if data[0] != '"' {
		newData := []byte{}
		newData = append(newData, '"')
		newData = append(newData, data...)
		newData = append(newData, '"')

		data = newData

		wasNumber = true
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errgo.Mask(err)
	}

	if err := parseDockerPort(s, d); err != nil {
		return errgo.Mask(err)
	}

	if wasNumber {
		d.formatJsonMode = modePortJsonNumber
	}

	return nil
}

// Empty returns true if this port is equal to "", false otherwise.
func (d *DockerPort) Empty() bool {
	return d.Port == "" // Protocol can be set automatically to TCP so don't check that
}

func (d *DockerPort) Equals(other DockerPort) bool {
	return d.Port == other.Port && d.Protocol == other.Protocol
}

func parseDockerPort(input string, dp *DockerPort) error {
	s := strings.Split(input, "/")

	switch len(s) {
	case 1:
		dp.Port = s[0]
		dp.Protocol = ProtocolTCP
		dp.formatJsonMode = modePortJsonString
	case 2:
		dp.Port = s[0]
		dp.Protocol = s[1]
		dp.formatJsonMode = modePortJsonDocker
	default:
		return errgo.Newf("Invalid format, must be either <port> or <port>/<prot>, got '%s'", input)
	}

	if parsedPort, err := strconv.Atoi(dp.Port); err != nil {
		return errgo.Notef(err, "Port must be a number, got '%s'", dp.Port)
	} else if parsedPort < 1 || parsedPort > 65535 {
		return errgo.Notef(err, "Port must be a number between 1 and 65535, got '%s'", dp.Port)
	}

	switch dp.Protocol {
	case "":
		return errgo.Newf("Protocol must not be empty.")
	case ProtocolUDP:
		fallthrough
	case ProtocolTCP:
		return nil
	default:
		return errgo.Newf("Unknown protocol: '%s' in '%s'", dp.Protocol, input)
	}
}
