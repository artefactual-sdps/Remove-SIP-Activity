package remove_test

import (
	"testing"

	activity "github.com/artefactual-sdps/remove-files-activity"
	temporalsdk_activity "go.temporal.io/sdk/activity"
	temporalsdk_testsuite "go.temporal.io/sdk/testsuite"
	"gotest.tools/v3/assert"
	tfs "gotest.tools/v3/fs"
)

func TestRemoveSipFiles(t *testing.T) {
	t.Parallel()

	td := tfs.NewDir(t, "remove-sip-files-test",
		tfs.WithDir(".DS_Store",
			tfs.WithFile("test", "hello from test"),
		),
		tfs.WithFile("keepme", "don't delete me."),
	)
	config := tfs.NewFile(t, ".ignore", tfs.WithContent(".DS_Store\n"))

	type Test struct {
		name    string
		params  activity.RemoveFilesParams
		want    *activity.RemoveFilesResult
		wantFs  tfs.Manifest
		wantErr string
	}
	for _, tt := range []Test{
		{
			name: "Should remove .DS_Store from directory",
			params: activity.RemoveFilesParams{
				RemovePath: td.Path(),
				IgnorePath: config.Path(),
			},
			want: &activity.RemoveFilesResult{Removed: []string{
				td.Join(
					".DS_Store",
				),
			}},
			wantFs: tfs.Expected(t, tfs.WithFile("keepme", "don't delete me.")),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := &temporalsdk_testsuite.WorkflowTestSuite{}
			env := ts.NewTestActivityEnvironment()
			env.RegisterActivityWithOptions(
				activity.NewRemoveFilesActivity().Execute,
				temporalsdk_activity.RegisterOptions{
					Name: activity.RemoveFilesName,
				},
			)

			enc, err := env.ExecuteActivity(activity.RemoveFilesName, tt.params)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)

				return
			}
			assert.NilError(t, err)

			var got activity.RemoveFilesResult
			if err = enc.Get(&got); err != nil {
				t.Fatalf("get results: %v", err)
			}

			assert.NilError(t, err)
			assert.DeepEqual(t, &got, tt.want)
			assert.Assert(t, tfs.Equal(td.Path(), tt.wantFs))
		})
	}
}
