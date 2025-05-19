package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

// Todo : Add test for model initialized with multiple directories
// Todo : Add test for clipboard different variations, cut paste
// Todo : Add test for tea resizing
// Todo : Add test for quitting

func TestCopy(t *testing.T) {
	curTestDir := filepath.Join(testDir, "TestCopy")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")
	t.Run("Basic Copy", func(t *testing.T) {
		setupDirectories(t, curTestDir, dir1, dir2)
		setupFiles(t, file1)
		t.Cleanup(func() {
			os.RemoveAll(curTestDir)
		})

		m := defaultTestModel(dir1)

		// Todo validate current panel is "dir1"
		// Todo : Move all basic validation to a separate test
		// Everything that doesn't have anything to do with copy paste

		// validate file1
		// Todo : improve the interface we use to interact with filepanel

		// Todo : file1.txt should not be duplicated

		// Todo : Having to send a random keypress to initiate model init.
		// Should not have to do that
		TeaUpdateWithErrCheck(t, &m, nil)

		assert.Equal(t, "file1.txt",
			m.fileModel.filePanels[m.filePanelFocusIndex].element[0].name)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.CopyItems[0]))

		// Todo : validate clipboard
		assert.False(t, m.copyItems.cut)
		assert.Equal(t, file1, m.copyItems.items[0])

		// move to dir2
		m.updateCurrentFilePanelDir("../dir2")
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Actual paste may take time, since its an os operations
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1.txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)

		// Todo : still on clipboard
		assert.False(t, m.copyItems.cut)
		assert.Equal(t, file1, m.copyItems.items[0])

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Actual paste may take time, since its an os operations
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1(1).txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)
		assert.FileExists(t, filepath.Join(dir2, "file1(1).txt"))
		// Todo : Also validate process bar having two processes.
	})
}
