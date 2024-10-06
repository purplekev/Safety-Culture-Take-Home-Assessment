package folder

import (
	"errors"
	"strings"
)

var (
	ErrInvalidSrcFolder = errors.New("error: Source folder does not exist")
	ErrInvalidDstFolder = errors.New("error: Destination folder does not exist")
	ErrMoveToSame       = errors.New("error: Cannot move a folder to itself")
	ErrMoveToChildOf    = errors.New("error: Cannot move a folder to a child of itself")
	ErrMoveToDiffOrg    = errors.New("error: Cannot move a folder to a different organization")
)

func (f *driver) MoveFolder(src string, dst string) ([]Folder, error) {
	all_folders := f.folders

	src_folder, exist_src := includesFolder(all_folders, src)
	dst_folder, exist_dst := includesFolder(all_folders, dst)

	if !exist_src { // invalid src
		return []Folder{}, ErrInvalidSrcFolder
	} else if !exist_dst { // invalid dst
		return []Folder{}, ErrInvalidDstFolder
	} else if src == dst { // moving to itself
		return []Folder{}, ErrMoveToSame
	} else if src_folder.OrgId != dst_folder.OrgId { // moving folder to diff orgID
		return []Folder{}, ErrMoveToDiffOrg
	} else if strings.Contains(dst_folder.Paths, src) { // moving to child of itself
		return []Folder{}, ErrMoveToChildOf
	}

	// deep copy of all_folders to modify
	res := make([]Folder, len(all_folders))
	copy(res, all_folders)

	new_path := dst_folder.Paths + "." + src_folder.Name
	for i := range res {
		if res[i].OrgId != src_folder.OrgId {
			continue
		} else if res[i].Name == src { // updating root folder
			res[i].Paths = new_path
		} else if strings.HasPrefix(res[i].Paths, src_folder.Paths+".") { // check for if its a child in one line
			res[i].Paths = strings.Replace(res[i].Paths, src_folder.Paths, new_path, 1) // update the child folder paths
		}
	}

	return res, nil
}
