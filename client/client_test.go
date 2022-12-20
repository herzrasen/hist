package client

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("test.db")
	require.NoError(t, err)
	err = c.Update("sudo dnf update")
	require.NoError(t, err)
}
