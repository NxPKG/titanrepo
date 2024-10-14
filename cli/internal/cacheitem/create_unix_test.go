//go:build darwin || linux
// +build darwin linux

package cacheitem

import (
	"syscall"
	"testing"

	"github.com/khulnasoft/titanrepo/cli/internal/titanpath"
	"gotest.tools/v3/assert"
)

func createFifo(t *testing.T, anchor titanpath.AbsoluteSystemPath, fileDefinition createFileDefinition) error {
	t.Helper()
	path := fileDefinition.Path.RestoreAnchor(anchor)
	fifoErr := syscall.Mknod(path.ToString(), syscall.S_IFIFO|0666, 0)
	assert.NilError(t, fifoErr, "FIFO")
	return fifoErr
}
