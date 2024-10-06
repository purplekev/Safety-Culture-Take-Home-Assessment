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
	ErrInvalidSrcFolder = errors.New("error: Source folder does not exist")
	ErrInvalidDstFolder = errors.New("error: Destination folder does not exist")
	ErrMoveToSame       = errors.New("error: Cannot move a folder to itself")
	ErrMoveToChildOf    = errors.New("error: Cannot move a folder to a child of itself")
	ErrMoveToDiffOrg    = errors.New("error: Cannot move a folder to a different organization")
)

var org1 = uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")
var org2 = uuid.FromStringOrNil("9b4cdb0a-cfea-4f9d-8a68-24f038fae385")

// function to retrieve sample data to test cases from spec
func getMoveFolderData() []folder.Folder {
	return []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: org1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: org1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: org1},
		{Name: "delta", Paths: "alpha.delta", OrgId: org1},
		{Name: "echo", Paths: "alpha.delta.echo", OrgId: org1},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: org2},
		{Name: "golf", Paths: "golf", OrgId: org1},
	}
}

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()

	// for these tests, i created a mini data set to test.
	// my thought process behind this, was that if we used the sample.json like for get_folder,
	// this would be looping through data set of 300+ entries unneccessarily for the same testing.
	tests := [...]struct {
		testName string
		src      string
		dst      string
		want     Result
	}{
		{
			testName: "Moving folder to connected tree with single child",
			src:      "bravo",
			dst:      "delta",
			want: Result{Folders: []folder.Folder{
				{Name: "alpha", OrgId: org1, Paths: "alpha"},
				{Name: "bravo", OrgId: org1, Paths: "alpha.delta.bravo"},
				{Name: "charlie", OrgId: org1, Paths: "alpha.delta.bravo.charlie"},
				{Name: "delta", OrgId: org1, Paths: "alpha.delta"},
				{Name: "echo", OrgId: org1, Paths: "alpha.delta.echo"},
				{Name: "foxtrot", OrgId: org2, Paths: "foxtrot"},
				{Name: "golf", OrgId: org1, Paths: "golf"},
			},
				Err: nil},
		},
		{
			testName: "Moving folder to unconnected folder",
			src:      "bravo",
			dst:      "golf",
			want: Result{Folders: []folder.Folder{
				{Name: "alpha", OrgId: org1, Paths: "alpha"},
				{Name: "bravo", OrgId: org1, Paths: "golf.bravo"},
				{Name: "charlie", OrgId: org1, Paths: "golf.bravo.charlie"},
				{Name: "delta", OrgId: org1, Paths: "alpha.delta"},
				{Name: "echo", OrgId: org1, Paths: "alpha.delta.echo"},
				{Name: "foxtrot", OrgId: org2, Paths: "foxtrot"},
				{Name: "golf", OrgId: org1, Paths: "golf"},
			},
				Err: nil},
		},
		{
			testName: "Moving folder with no children",
			src:      "charlie",
			dst:      "echo",
			want: Result{Folders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: org1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: org1},
				{Name: "charlie", Paths: "alpha.delta.echo.charlie", OrgId: org1},
				{Name: "delta", Paths: "alpha.delta", OrgId: org1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: org1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: org2},
				{Name: "golf", Paths: "golf", OrgId: org1},
			},
				Err: nil},
		},
		{
			testName: "Moving entire folder tree to single folder",
			src:      "alpha",
			dst:      "golf",
			want: Result{Folders: []folder.Folder{
				{Name: "alpha", Paths: "golf.alpha", OrgId: org1},
				{Name: "bravo", Paths: "golf.alpha.bravo", OrgId: org1},
				{Name: "charlie", Paths: "golf.alpha.bravo.charlie", OrgId: org1},
				{Name: "delta", Paths: "golf.alpha.delta", OrgId: org1},
				{Name: "echo", Paths: "golf.alpha.delta.echo", OrgId: org1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: org2},
				{Name: "golf", Paths: "golf", OrgId: org1},
			},
				Err: nil},
		},
		{
			testName: "Err: Source folder does not exist",
			src:      "invalid_folder",
			dst:      "delta",
			want:     Result{Folders: []folder.Folder{}, Err: ErrInvalidSrcFolder},
		},
		{
			testName: "Err: Destination folder does not exist",
			src:      "bravo",
			dst:      "invalid_folder",
			want:     Result{Folders: []folder.Folder{}, Err: ErrInvalidDstFolder},
		},
		{
			testName: "Err: Cannot move a folder to itself",
			src:      "bravo",
			dst:      "bravo",
			want:     Result{Folders: []folder.Folder{}, Err: ErrMoveToSame},
		},
		{
			testName: "Err: Cannot move a folder to a child of itself",
			src:      "bravo",
			dst:      "charlie",
			want:     Result{Folders: []folder.Folder{}, Err: ErrMoveToChildOf},
		},
		{
			testName: "Err: Cannot move a folder to a different organization",
			src:      "bravo",
			dst:      "foxtrot",
			want:     Result{Folders: []folder.Folder{}, Err: ErrMoveToDiffOrg},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			f := folder.NewDriver(getMoveFolderData())
			getNewStructure, getErr := f.MoveFolder(tt.src, tt.dst)
			assert.Equal(t, tt.want.Folders, getNewStructure)
			assert.Equal(t, tt.want.Err, getErr)
		})
	}
}
