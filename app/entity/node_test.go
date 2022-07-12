package entity

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var root Node
var data []Node

//+--~/ (id=0)
//|  +--tung/ (id=1)
//|  |  +--file1/ (id=2)
//|  |  +--file2/ (id=3)
//|  |  +--folder1/ (id=4)
//|  |  |  +--file3/ (id=5)
//|  +--usr/ (id=6)
//|  |  +--file4/ (id=7)

func initTestData() {
	root = Node{
		Id:       "0",
		Type:     0,
		Name:     "~",
		Data:     "",
		Path:     "/",
		Children: []*Node{},
	}
	data = append(data, root)

	x1 := Node{
		Id:       "1",
		Type:     0,
		Parent:   &root,
		Name:     "tung",
		Data:     "",
		Path:     "/tung",
		Children: []*Node{},
		Root:     &root,
	}
	root.Children = append(root.Children, &x1)
	data = append(data, x1)

	x2a := Node{
		Id:       "2",
		Type:     1,
		Parent:   &x1,
		Name:     "file1",
		Path:     "/tung/file1",
		Data:     "hello world",
		Children: []*Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2a)
	data = append(data, x2a)

	x2b := Node{
		Id:       "3",
		Type:     1,
		Parent:   &x1,
		Name:     "file2",
		Path:     "/tung/file2",
		Data:     "as98d9asd",
		Children: []*Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2b)
	data = append(data, x2b)

	x2c := Node{
		Id:       "4",
		Type:     0,
		Parent:   &x1,
		Name:     "folder1",
		Path:     "/tung/folder1",
		Data:     "",
		Children: []*Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2c)
	data = append(data, x2c)

	x2d := Node{
		Id:       "5",
		Type:     1,
		Parent:   &x2c,
		Name:     "file3",
		Path:     "/tung/folder1/file3",
		Data:     "as98d9asd",
		Children: []*Node{},
		Root:     &root,
	}
	x2c.Children = append(x2c.Children, &x2d)
	data = append(data, x2d)

	x3 := Node{
		Id:       "6",
		Type:     0,
		Parent:   &root,
		Name:     "usr",
		Path:     "/usr",
		Data:     "",
		Children: []*Node{},
		Root:     &root,
	}
	root.Children = append(root.Children, &x3)
	data = append(data, x3)

	x4 := Node{
		Id:       "7",
		Type:     1,
		Parent:   &x3,
		Name:     "file4",
		Path:     "/usr/file4",
		Data:     "hello tung",
		Children: []*Node{},
		Root:     &root,
	}
	x3.Children = append(x3.Children, &x4)
	data = append(data, x4)

}

func TestParsePathToNodeNameSlice(t *testing.T) {
	path := "tung/hello/world"
	want := []string{"tung", "hello", "world"}
	got, err := ParsePathToItemViewNameSlice(path)
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	path = "tung/hello/world/"
	want = []string{"tung", "hello", "world"}
	got, err = ParsePathToItemViewNameSlice(path)
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	path = "/tung/hello/world"
	want = []string{"tung", "hello", "world"}
	got, err = ParsePathToItemViewNameSlice(path)
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	path = "tung/hello/wo rld"
	want = []string{"tung", "hello", "wo rld"}
	got, err = ParsePathToItemViewNameSlice(path)
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	path = "./tung"
	want = []string{"tung"}
	got, err = ParsePathToItemViewNameSlice(path)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNode_TraverseByPath(t *testing.T) {
	initTestData()
	root.Print(0, "")

	t.Run("normal cases", func(t *testing.T) {
		gotFolder, err := root.TraverseByPath("/tung")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)

		gotFolder, err = root.TraverseByPath("/tung/")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)

		gotFolder, err = root.TraverseByPath("/usr/file3")
		assert.Equal(t, err, ErrorInvalidPath)
		//
		// to root
		gotFolder, err = root.TraverseByPath("/")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())

		// current folder
		gotFolder, err = root.TraverseByPath(".")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())
		//
		gotFolder, err = root.TraverseByPath("./")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())

		// parent folder of root
		gotFolder, err = root.TraverseByPath("..")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())
		//
		gotFolder, err = root.TraverseByPath("../")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())

		// parent folder
		node := root.Children[0]
		gotFolder, err = node.TraverseByPath("..")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())
		//
		gotFolder, err = node.TraverseByPath("../")
		assert.NoError(t, err)
		assert.Equal(t, true, gotFolder.IsRoot())
	})

	t.Run("navigate using current folder alias", func(t *testing.T) {
		// possible cases:
		// - ./path/to
		gotFolder, err := root.TraverseByPath("./tung")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)

		gotFolder, err = root.TraverseByPath("./tung/")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)

		gotFolder, err = gotFolder.TraverseByPath("./")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)

		gotFolder, err = gotFolder.TraverseByPath(".")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)
	})

	t.Run("navigate using parent folder alias", func(t *testing.T) {
		// possible cases:
		// - ../path/to
		// - ../path/../to
		var gotFolder *Node
		var err error

		gotFolder, err = root.TraverseByPath("..")
		assert.NoError(t, err)
		assert.Equal(t, "~", gotFolder.Name)

		gotFolder, err = root.TraverseByPath("../")
		assert.NoError(t, err)
		assert.Equal(t, "~", gotFolder.Name)

		gotFolder, err = root.TraverseByPath("../..")
		assert.NoError(t, err)
		assert.Equal(t, "~", gotFolder.Name)

		gotFolder, err = root.TraverseByPath("/tung/..")
		assert.NoError(t, err)
		assert.Equal(t, "~", gotFolder.Name)
		//
		gotFolder, err = root.TraverseByPath("/tung/folder1/..")
		assert.NoError(t, err)
		assert.Equal(t, "tung", gotFolder.Name)
		//
		gotFolder, err = root.TraverseByPath("/tung/folder1/../..")
		assert.NoError(t, err)
		assert.Equal(t, "~", gotFolder.Name)
		//
		//// wrong cases
		gotFolder, err = root.TraverseByPath("/tung/usr/../..")
		assert.Equal(t, ErrorInvalidPath, err)

		gotFolder, err = root.TraverseByPath("/tung/file1.txt/..")
		assert.Equal(t, ErrorInvalidPath, err)
	})

	t.Run("navigate using home alias", func(t *testing.T) {
		// possible cases:
		// - when at a folder, navigate to any folder using '/' alias
		var folder1 *Node
		var gotFolder *Node
		var err error
		log.Println(gotFolder)

		folder1, err = root.TraverseByPath("./tung")
		assert.NoError(t, err)
		assert.Equal(t, "tung", folder1.Name)

		// at /tung, navigate to root using '/'
		gotFolder, err = folder1.TraverseByPath("/")
		assert.NoError(t, err)
		assert.Equal(t, "~", gotFolder.Name)

		// at /tung, navigate to /usr using '/'
		gotFolder, err = folder1.TraverseByPath("/usr")
		assert.NoError(t, err)
		assert.Equal(t, "usr", gotFolder.Name)
	})

	t.Run("invalid path", func(t *testing.T) {
		var err error

		_, err = root.TraverseByPath("/.")
		assert.Equal(t, ErrorInvalidPath, err)

		_, err = root.TraverseByPath("/tun/..")
		assert.Equal(t, ErrorInvalidPath, err)

		_, err = root.TraverseByPath("/%/..")
		assert.Equal(t, ErrorInvalidPath, err)

		_, err = root.TraverseByPath("///..")
		assert.Equal(t, ErrorInvalidPath, err)
	})
}

func TestNode_InsertByPath(t *testing.T) {
	initTestData()
	root.Print(0, "")
	var err error

	t.Run("normal", func(t *testing.T) {

		insertedNode := Node{
			Type:     1,
			Parent:   &data[1],
			Name:     "new_file",
			Path:     "",
			Data:     "",
			Children: nil,
		}

		err := root.InsertByPath("/tung/new_file", insertedNode)
		assert.NoError(t, err)
		root.Print(0, "")

		insertedNode2 := Node{
			Type:     0,
			Parent:   nil,
			Name:     "newfolder",
			Path:     "",
			Data:     "",
			Children: nil,
		}

		err = root.InsertByPath("/tung/../newfolder", insertedNode2)
		assert.NoError(t, err)
		gotNode, err := root.TraverseByPath("/newfolder")
		assert.NoError(t, err)
		assert.Equal(t, "newfolder", gotNode.Name)

		insertedNode3 := Node{
			Type:     0,
			Parent:   nil,
			Name:     "newfolder2",
			Path:     "",
			Data:     "",
			Children: nil,
		}

		err = root.InsertByPath("./newfolder2", insertedNode3)
		assert.NoError(t, err)
		gotNode, err = root.TraverseByPath("/newfolder2")
		assert.NoError(t, err)
		assert.Equal(t, "newfolder2", gotNode.Name)
	})

	t.Run("wrong path case", func(t *testing.T) {
		insertedNode2 := Node{
			Type:     0,
			Parent:   nil,
			Name:     "newfolder",
			Path:     "",
			Data:     "",
			Children: nil,
		}

		err = root.InsertByPath("/tug/newfolder", insertedNode2)
		assert.Equal(t, err, ErrorInvalidPath)

		err = root.InsertByPath("/./newfolder", insertedNode2)
		assert.Equal(t, err, ErrorInvalidPath)

		err = root.InsertByPath("//tung/newfolder", insertedNode2)
		assert.Equal(t, err, ErrorInvalidPath)

		err = root.InsertByPath("/tung/newfolder/newfolder", insertedNode2)
		assert.Equal(t, err, ErrorInvalidPath)

		err = root.InsertByPath("/tung/../usr/fake_folder/newfolder", insertedNode2)
		assert.Equal(t, err, ErrorInvalidPath)
	})
}

func TestNode_FindByName(t *testing.T) {
	initTestData()
	root.Print(0, "")

	gotNodes := root.FindByName("folder1")
	assert.Equal(t, 1, len(gotNodes))

	// fail case
	gotNodes = root.FindByName("file2312")
	assert.Equal(t, 0, len(gotNodes))

	//// substring case
	gotNodes = root.FindByName("file")
	assert.Equal(t, 4, len(gotNodes))
}

func TestEditNodeValues(t *testing.T) {
	initTestData()
	root.Print(0, "")

	root.Children[0].Name = "new name here"
	root.Print(0, "")
}

func TestNode_DeleteByPath(t *testing.T) {
	initTestData()
	root.Print(0, "")

	t.Run("normal", func(t *testing.T) {
		gotErr := root.DeleteByPath("/tung/file1")
		assert.NoError(t, gotErr)

		root.Print(0, "")

		gotNodes := root.FindByName("file1")
		assert.Equal(t, 0, len(gotNodes))
	})

	t.Run("normal", func(t *testing.T) {
		gotErr := root.DeleteByPath("/usr")
		assert.NoError(t, gotErr)

		root.Print(0, "")

		gotNodes := root.FindByName("usr")
		assert.Equal(t, 0, len(gotNodes))
	})

	t.Run("folder case", func(t *testing.T) {
		gotErr := root.DeleteByPath("/tung/folder1")
		assert.NoError(t, gotErr)

		gotNodes := root.FindByName("folder1")
		assert.Equal(t, 0, len(gotNodes))
	})

	t.Run("wrong name case", func(t *testing.T) {
		gotErr := root.DeleteByPath("/tung/file_wrong.txt")
		assert.Equal(t, ErrorInvalidPath, gotErr)

	})

	t.Run("delete root", func(t *testing.T) {
		gotErr := root.DeleteByPath("/")
		assert.Equal(t, ErrorCantDeleteRoot, gotErr)

		root.Print(0, "")
	})

	t.Run("delete current folder", func(t *testing.T) {
		node := data[1]
		gotErr := node.DeleteByPath(".")
		assert.Equal(t, ErrorCantDeleteCurrentNodeOrParent, gotErr)

		gotErr = node.DeleteByPath("./")
		assert.Equal(t, ErrorCantDeleteCurrentNodeOrParent, gotErr)

		root.Print(0, "")
	})

	t.Run("delete current folder", func(t *testing.T) {
		node := data[1]
		gotErr := node.DeleteByPath("..")
		assert.Equal(t, ErrorCantDeleteCurrentNodeOrParent, gotErr)

		gotErr = node.DeleteByPath("../")
		assert.Equal(t, ErrorCantDeleteCurrentNodeOrParent, gotErr)

		root.Print(0, "")
	})

}

func TestNode_UpdatePreOrderNumber(t *testing.T) {
	initTestData()
	root.Print(0, "")

	root.UpdatePreOrderNumber()
	//
	root.PrintOrder(0, "")
	//
	//fmt.Printf("total node: %d\n", root.TotalNodes())
}
