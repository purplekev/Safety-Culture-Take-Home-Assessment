package folder

import (
	"strings"
)

// moves a subtree from one parent folder to another maintaining the order of the children.
// returns new folder structure or an error otherwise.
func (f *driver) MoveFolder(src string, dst string) ([]Folder, error) {
	all_folders := f.folders

	src_folder, exist_src := includesFolder(all_folders, src)
	dst_folder, exist_dst := includesFolder(all_folders, dst)

	// error checks
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

	// copy of all_folders to modify
	res := make([]Folder, len(all_folders))
	copy(res, all_folders)

	new_path := dst_folder.Paths + "." + src_folder.Name
	for i := range res {
		if res[i].OrgId != src_folder.OrgId { // cannot manipulate folders if they are in different
			continue
		} else if res[i].Name == src { // updating parent folder
			res[i].Paths = new_path
		} else if strings.HasPrefix(res[i].Paths, src_folder.Paths+".") { // updating child nodes affected
			res[i].Paths = strings.Replace(res[i].Paths, src_folder.Paths, new_path, 1)
		}
	}

	return res, nil
}
