package lib

import (
	"regexp"
	"strings"
)

func InterfaceToString(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func CheckName(name string) bool {
	if len(name) == 0 {
		return false
	}
	ok, _ := regexp.MatchString(NAME_REGEX, name)
	if !ok {
		return false
	}
	return true
}

const (
	Folder     int = 0
	File           = 1
	NAME_REGEX     = "^[a-zA-Z0-9 _-]+$"
)

//CleanPathToFileOrFolder remove redundant slash of path (if any)
//
// Return same path for alias ('/', '.', '..') and valid path (/path/to)
//
// Treat alias as a valid folder name
func CleanPathToFileOrFolder(path string) string {
	if !strings.HasPrefix(path, ".") && !strings.HasPrefix(path, "/") {
		return ""
	}
	if IsCurrentFolder(path) || IsParentFolder(path) || path == "/" {
		return path
	}

	// remove last slash
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	tmp := make([]string, len(path))
	split := strings.Split(path, "")
	copy(tmp, split)
	//fmt.Printf("tmp: %v, split: %v\n", tmp, split)
	copy := strings.Join(tmp, "")

	// remove first slash for processing
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	parts := strings.Split(path, "/")
	for i, x := range parts {
		//fmt.Printf("check x = %s\n", x)
		if x == ".." {
			continue
		} else if x == "." && i == 0 {
			if copy[0] == '/' {
				return ""
			}
			continue
		} else if x == "." && i != 0 {
			return ""
		} else if CheckName(x) == false {
			return ""
		}
	}

	return copy
}

//GetParentPath return the path of parent of the last item
//
// path must be a valid
func GetParentPath(path string) string {
	if path == "/" {
		return path
	} else if IsCurrentFolder(path) || IsParentFolder(path) {
		return ""
	}
	nodeNames := strings.Split(path, "/")
	if len(nodeNames) == 2 {
		return nodeNames[0] + "/"
	}

	return strings.Join(nodeNames[:len(nodeNames)-1], "/")
}

func IsCurrentFolder(path string) bool {
	if path == "./" || path == "." {
		return true
	}
	return false
}

func IsParentFolder(path string) bool {
	if path == "../" || path == ".." {
		return true
	}
	return false
}
