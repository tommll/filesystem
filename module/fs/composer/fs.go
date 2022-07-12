package composer

import "fs/app/entity"

type DirectoryContentDisplay struct {
	Items []string `json:"items"`
}

func NodeListToFileList(list []*entity.Node) DirectoryContentDisplay {
	res := DirectoryContentDisplay{Items: []string{}}
	for _, x := range list {
		res.Items = append(res.Items, x.Name)
	}
	return res
}
