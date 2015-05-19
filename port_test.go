package generictypes

import (
	"encoding/json"
	"testing"
)

var validPorts = []struct {
	Input string
}{
	{"1/tcp"},
	{"65000/tcp"},
	{"23/udp"},
	{"65000/udp"},
}

func TestDockerPort__ValidPorts(t *testing.T) {
	for _, data := range validPorts {
		_, err := ParseDockerPort(data.Input)
		if err != nil {
			t.Fatalf("Expected no error for input: %v\nBut got: %#v", data.Input, err)
		}
	}
}

var invalidPorts = []struct {
	Input string
}{
	{""},        // empty
	{"a/tcp"},   // wrong port
	{"90/icmp"}, // wrong protocol
	{"a/b/c"},   // too many slashes

	{"0/tcp"},
	{"66000/udp"},

	{"/"},
	{"/80/"},
	{"/80"},
	{"/tcp"},
	{"tcp/"},
	{"invalid/"},
	{"/invalid"},
	{"-80/tcp"},
	{"-80"},
}

func TestDockerPortParsingErrors(t *testing.T) {
	for _, data := range invalidPorts {
		port, err := ParseDockerPort(data.Input)
		if err == nil {
			t.Fatalf("Expected error for input: %v\nBut got: %#v", data.Input, port)
		}
	}
}

var validDockerPortJsonInput = []struct {
	Input    string
	Port     string
	Protocol string
}{
	{"80", "80", ProtocolTCP},               // port as int
	{"\"8080\"", "8080", ProtocolTCP},       // port as string
	{"\"10000/udp\"", "10000", ProtocolUDP}, // port and protocol in docker notation
}

func TestDockerPort__ValidJSONInput(t *testing.T) {
	for _, data := range validDockerPortJsonInput {
		var port DockerPort
		if err := json.Unmarshal([]byte(data.Input), &port); err != nil {
			t.Fatalf("Expected no error for input: %v\nBut got: %#v", data.Input, err)
		}

		if port.Port != data.Port {
			t.Fatalf("Expected '%s' but got '%s' as Port", data.Port, port.Port)
		}
		if port.Protocol != data.Protocol {
			t.Fatalf("Expected '%s' but got '%s' as Port", data.Protocol, port.Protocol)
		}
	}
}

var invalidDockerPortJsonInput = []struct {
	Input string
}{
	{"[]"},
	{"{]"},
	{"asdf"},
	{""},
	{"{}"},
}

func TestDockerPort__InalidJSONInput(t *testing.T) {
	for _, data := range invalidDockerPortJsonInput {
		var port DockerPort
		if err := json.Unmarshal([]byte(data.Input), &port); err == nil {
			t.Fatalf("Expected error for input: %v\nBut got: %#v", data.Input, port)
		}

	}
}
