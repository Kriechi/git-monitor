package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCmd(t *testing.T) {
	buf := new(bytes.Buffer)

	root := newRootCmd(buf)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"list", "-v", "--repo_dir", "/tmp"})

	c, err := root.ExecuteC()
	assert.NotNil(t, c)
	assert.Nil(t, err)

	result := buf.String()
	assert.Contains(t, result, "REPOSITORY")
	assert.Contains(t, result, "REMOTE URL")
	assert.Contains(t, result, "MONITORED BRANCHES")
}
