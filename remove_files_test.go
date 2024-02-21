package main

import (
	"testing"

	temporalsdk_activity "go.temporal.io/sdk/activity"
	temporalsdk_testsuite "go.temporal.io/sdk/testsuite"
	"gotest.tools/v3/assert"
	tfs "gotest.tools/v3/fs"
)

func TestRemoveSipFiles(t *testing.T) {

	t.Parallel()

	td := tfs.NewDir(t, "remove-sip-files-test", tfs.WithDir(".DS_Store", tfs.WithFile("test", "hello from test")))
	config := tfs.NewFile(t, ".succumb", tfs.WithContent(`
			.DS_Store/
		`))

	type Test struct {
		name    string
		params  RemoveSIPFilesParams
		wantErr string
	}
	for _, tt := range []Test{
		{
			name: "Should remove ds store from directory",
			params: RemoveSIPFilesParams{
				SipPath:     td.Path(),
				SuccumbPath: config.Path(),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := &temporalsdk_testsuite.WorkflowTestSuite{}
			env := ts.NewTestActivityEnvironment()
			env.RegisterActivityWithOptions(
				NewRemoveSIPFilesActivity().Execute,
				temporalsdk_activity.RegisterOptions{
					Name: RemoveSIPFilesName,
				},
			)
			_, err := env.ExecuteActivity(RemoveSIPFilesName, tt.params)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
			}

			assert.NilError(t, err)

		})
	}

}
