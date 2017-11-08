package pathutil

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChangeDirForFunction(t *testing.T) {
	origDir, err := CurrentWorkingDirectoryAbsolutePath()
	require.NoError(t, err)

	// now change dir, but just for the function
	newDir := UserHomeDir()
	require.NoError(t, ChangeDirForFunction(newDir, func() {
		// current dir should be the changed value
		dir, err := CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)
		require.Equal(t, newDir, dir)
	}))

	// current dir should be the original value
	dir, err := CurrentWorkingDirectoryAbsolutePath()
	require.NoError(t, err)
	require.Equal(t, origDir, dir)
}

func TestRevokableChangeDir(t *testing.T) {
	origDir, err := CurrentWorkingDirectoryAbsolutePath()
	require.NoError(t, err)

	// revokable change dir
	newDir := UserHomeDir()
	revokeFn, err := RevokableChangeDir(newDir)
	require.NoError(t, err)

	{
		// current dir should be the changed value
		dir, err := CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)
		require.Equal(t, newDir, dir)
	}

	{
		// revoke it
		require.NoError(t, revokeFn())

		// current dir should be the original value
		dir, err := CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)
		require.Equal(t, origDir, dir)
	}

}

func TestEnsureDirExist(t *testing.T) {
	// testDir not exist

	currDirPath, err := filepath.Abs(".")
	require.Equal(t, nil, err)

	currentTime := time.Now()
	currentTimeStamp := currentTime.Format("20060102150405")
	testDir := path.Join(currDirPath, currentTimeStamp+"TestEnsurePathExist")
	exist, err := IsDirExists(testDir)
	require.Equal(t, nil, err)
	require.Equal(t, false, exist)
	defer func() {
		require.Equal(t, nil, os.Remove(testDir))
	}()

	require.Equal(t, nil, EnsureDirExist(testDir))
	exist, err = IsDirExists(testDir)
	require.Equal(t, nil, err)
	require.Equal(t, true, exist)

	// testDir exist

	require.Equal(t, nil, EnsureDirExist(testDir))
	exist, err = IsDirExists(testDir)
	require.Equal(t, nil, err)
	require.Equal(t, true, exist)
}

func TestIsRelativePath(t *testing.T) {
	// should return true if relative path, false if absolute path

	require.Equal(t, true, IsRelativePath("./rel"))
	require.Equal(t, false, IsRelativePath("/abs"))
	require.Equal(t, false, IsRelativePath("$THISENVDOESNTEXIST/some"))
	require.Equal(t, true, IsRelativePath("rel"))
}

func TestIsPathExists(t *testing.T) {
	// should return false if path doesn't exist

	exists, err := IsPathExists("this/should/not/exist")
	require.Equal(t, nil, err)
	require.Equal(t, false, exists)

	exists, err = IsPathExists(".")
	require.Equal(t, nil, err)
	require.Equal(t, true, exists)

	exists, err = IsPathExists("")
	require.NotEqual(t, nil, err)
	require.Equal(t, false, exists)
}

func TestAbsPath(t *testing.T) {
	// should expand path

	currDirPath, err := filepath.Abs(".")
	require.Equal(t, nil, err)
	require.NotEqual(t, "", currDirPath)
	require.NotEqual(t, ".", currDirPath)

	homePathEnv := "/path/home/test-user"
	require.Equal(t, nil, os.Setenv("HOME", homePathEnv))

	testFileRelPathFromHome := "some/file.ext"
	absPathToTestFile := fmt.Sprintf("%s/%s", homePathEnv, testFileRelPathFromHome)

	expandedPath, err := AbsPath("")
	require.NotEqual(t, nil, err)
	require.Equal(t, "", expandedPath)

	expandedPath, err = AbsPath(".")
	require.Equal(t, nil, err)
	require.Equal(t, currDirPath, expandedPath)

	expandedPath, err = AbsPath(fmt.Sprintf("$HOME/%s", testFileRelPathFromHome))
	require.Equal(t, nil, err)
	require.Equal(t, absPathToTestFile, expandedPath)

	expandedPath, err = AbsPath(fmt.Sprintf("~/%s", testFileRelPathFromHome))
	require.Equal(t, nil, err)
	require.Equal(t, absPathToTestFile, expandedPath)
}

func TestUserHomeDir(t *testing.T) {
	// should return the path of the users home directory

	require.NotEqual(t, "", UserHomeDir())
}

func TestNormalizedOSTempDirPath(t *testing.T) {
	// returned temp dir path should not have a / at it's end

	tmpPth, err := NormalizedOSTempDirPath("some-test")
	require.Equal(t, nil, err)
	require.Equal(t, false, strings.HasSuffix(tmpPth, "/"))

	// should work if empty prefix is defined
	tmpPth, err = NormalizedOSTempDirPath("")
	require.Equal(t, nil, err)
	require.Equal(t, false, strings.HasSuffix(tmpPth, "/"))
}