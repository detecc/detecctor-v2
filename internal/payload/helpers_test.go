package payload

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseIP(t *testing.T) {
	require := require.New(t)

	ip := "192.168.1.1:8090"
	parsedIp, port := ParseIP(ip)

	require.Equal("192.168.1.1", parsedIp)
	require.Equal("8090", port)

	ip = "192.168.1.27"
	parsedIp, port = ParseIP(ip)

	require.Equal("", parsedIp)
	require.Equal("", port)
}
