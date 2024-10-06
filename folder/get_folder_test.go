package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// Pre defined variables for error handling
var (
	ErrInvalidOrg       = errors.New("error: OrgId does not exist or is invalid")
	ErrInvalidFolder    = errors.New("error: Folder does not exist")
	ErrFolderInvalidOrg = errors.New("error: Folder does not exist in the given org")
)

// function has return type result - either a slice of folders or an error
type Result struct {
	Folders []folder.Folder
	Err     error
}

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	sample_data := folder.GetSampleData()
	defaultOrgId := uuid.FromStringOrNil(folder.DefaultOrgID)

	tests := [...]struct {
		testName string
		orgID    uuid.UUID
		folders  []folder.Folder
		want     []folder.Folder
	}{
		{
			testName: "get folders with defaultOrgId",
			orgID:    defaultOrgId,
			folders:  sample_data,
			want:     sample_data[79:220],
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			get := f.GetFoldersByOrgID(tt.orgID)
			assert.Equal(t, tt.want, get)
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	// extract sample.json data
	sample_data := folder.GetSampleData()
	defaultOrgId := uuid.FromStringOrNil(folder.DefaultOrgID)
	var invalidOrgID uuid.UUID

	// data for specific edge case, parent and child have the same name
	same_name_parent_child_data := []folder.Folder{
		{
			Name:  "bob",
			OrgId: defaultOrgId,
			Paths: "bob",
		},
		{
			Name:  "bob",
			OrgId: defaultOrgId,
			Paths: "bob.bob",
		},
	}

	// expected output for specific edge case, parent and child have the same name
	same_name_parent_child_res := []folder.Folder{
		{
			Name:  "bob",
			OrgId: defaultOrgId,
			Paths: "bob.bob",
		},
	}

	tests := [...]struct {
		testName   string
		folderName string
		orgID      uuid.UUID
		folders    []folder.Folder // what we feed in
		want       Result          // expected output
	}{
		{
			testName:   "Multiple, multiple level hierachy",
			folderName: "settling-hobgoblin",
			orgID:      defaultOrgId,
			folders:    sample_data,
			want:       Result{Folders: sample_data[132:139], Err: nil},
		},
		{
			testName:   "Multiple, single level hierachy",
			folderName: "dashing-forearm",
			orgID:      defaultOrgId,
			folders:    sample_data,
			want:       Result{Folders: sample_data[212:215], Err: nil},
		},
		{
			testName:   "Root folder - no children",
			folderName: "endless-azrael",
			orgID:      defaultOrgId,
			folders:    sample_data,
			want:       Result{Folders: []folder.Folder{}, Err: nil},
		},
		{
			testName:   "Single child",
			folderName: "super-cobweb",
			orgID:      defaultOrgId,
			folders:    sample_data,
			want:       Result{Folders: sample_data[130:131], Err: nil},
		},
		{
			testName:   "Parent and child have the same name eg. bob.bob",
			folderName: "bob",
			orgID:      defaultOrgId,
			folders:    same_name_parent_child_data,
			want:       Result{Folders: same_name_parent_child_res, Err: nil},
		},
		{
			testName:   "Err: Folder does not exist in specified organization",
			folderName: "merry-fantomex",
			orgID:      defaultOrgId,
			folders:    sample_data,
			want:       Result{Folders: []folder.Folder{}, Err: ErrFolderInvalidOrg},
		},
		{
			testName:   "Err: Folder does not exist",
			folderName: "invalid_folder",
			orgID:      defaultOrgId,
			folders:    sample_data,
			want:       Result{Folders: []folder.Folder{}, Err: ErrInvalidFolder},
		},
		{
			testName:   "Err: Invalid OrgId",
			folderName: "endless-azrael",
			orgID:      invalidOrgID,
			folders:    sample_data,
			want:       Result{Folders: []folder.Folder{}, Err: ErrInvalidOrg},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			getChildren, getErr := f.GetAllChildFolders(tt.orgID, tt.folderName)
			assert.Equal(t, tt.want.Folders, getChildren)
			assert.Equal(t, tt.want.Err, getErr)
		})
	}
}
