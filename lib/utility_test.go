package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckName(t *testing.T) {
	assert.Equal(t, false, CheckName(""))

	assert.Equal(t, true, CheckName("tung"))

	assert.Equal(t, true, CheckName("tun g 1"))

	assert.Equal(t, true, CheckName("d_a - 2"))
	//
	assert.Equal(t, false, CheckName("/asd"))

	assert.Equal(t, false, CheckName("/asd% +"))

	assert.Equal(t, false, CheckName("/asd/"))
}

func TestCleanPathToFileOrFolder(t *testing.T) {
	var path string
	var gotPath string

	path = "/"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "/path/to"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)
	//
	path = "/path/to/"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, "/path/to", gotPath)
	//
	path = "."
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)
	//
	path = "./"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "./tung"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "./tung/"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, "./tung", gotPath)

	path = ".."
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "../"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "../../path/to"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)
	//
	path = "../../path/to/"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, "../../path/to", gotPath)
	//
	path = "/path/.."
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "/path/../to"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)
	//
	path = "/path/../to/"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, "/path/../to", gotPath)
	//
	path = "/path/../to/.."
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, path, gotPath)

	path = "/path/../to/../"
	gotPath = CleanPathToFileOrFolder(path)
	assert.Equal(t, "/path/../to/..", gotPath)

	t.Run("wrong cases", func(t *testing.T) {
		path = "random string"
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "/..."
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = ".../"
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "/path."
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "/pa*/.."
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)
		//
		path = "../."
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "/."
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "///."
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "//"
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)

		path = "....123"
		gotPath = CleanPathToFileOrFolder(path)
		assert.Equal(t, "", gotPath)
	})
}

func TestGetParentPath(t *testing.T) {
	var path string
	//var gotPath string

	data := map[string]string{
		"0":  "/",
		"1":  "/path",
		"2":  "/path/to",
		"3":  ".",
		"4":  "./",
		"5":  "./tung",
		"6":  "..",
		"7":  "../",
		"8":  "../../path/to",
		"9":  "/path/..",
		"10": "../path"}

	path = data["0"]
	assert.Equal(t, "/", GetParentPath(path))

	path = data["1"]
	assert.Equal(t, "/", GetParentPath(path))

	path = data["2"]
	assert.Equal(t, "/path", GetParentPath(path))

	path = data["3"]
	assert.Equal(t, "", GetParentPath(path))

	path = data["4"]
	assert.Equal(t, "", GetParentPath(path))

	path = data["5"]
	assert.Equal(t, "./", GetParentPath(path))

	path = data["6"]
	assert.Equal(t, "", GetParentPath(path))

	path = data["7"]
	assert.Equal(t, "", GetParentPath(path))

	path = data["8"]
	assert.Equal(t, "../../path", GetParentPath(path))

	path = data["9"]
	assert.Equal(t, "/path", GetParentPath(path))

	path = data["10"]
	assert.Equal(t, "../", GetParentPath(path))
}
