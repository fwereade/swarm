package bundle

import (
	"testing"

	"github.com/mrcrowl/swarm/config"
	"github.com/mrcrowl/swarm/source"
	"github.com/mrcrowl/swarm/testutil"

	"github.com/stretchr/testify/assert"
)

const buildDescrSampleJSON = `{
	"modules": [
		{
			"name": "abcd/efgh",
			"include": [
				"common/util",
				"common/dict",
				"common/strings"
			]
		},
		{
			"name": "wxyz/zzzz",
			"exclude": [
				"abcd/efgh"
			]
		},
		{
			"name": "stuv/vvvv",
			"exclude": [
				"abcd/efgh",
				"wxyz/zzzz"
			]
		}
	],
    "base": "app/src/"
}`

// duplicated in other places
func createWorkspace() *source.Workspace {
	ws := source.NewWorkspace("c:\\wf\\lp\\web\\App")
	return ws
}

func TestCreateModuleSet(t *testing.T) {
	descr, err := config.LoadBuildDescriptionString(buildDescrSampleJSON)
	assert.Nil(t, err)

	assert.Len(t, descr.Modules, 3)
	workspacePath := testutil.CreateTempDir()
	testutil.WriteTextFile(workspacePath, "Config.js", "")
	set := CreateModuleSet(source.NewWorkspace(workspacePath), descr.NormaliseModules(workspacePath), config.NewRuntimeConfig("", ""))
	assert.True(t, assert.ObjectsAreEqual([]string{"abcd/efgh", "wxyz/zzzz", "stuv/vvvv"}, set.names()), "Module order doesn't match")
}

// func TestCreateModuleSetFromFile(t *testing.T) {
// 	descr, err := config.LoadBuildDescriptionFile("c:\\wf\\lp\\web\\App\\build\\systemjs_build_controlpanel.json")
// 	assert.Nil(t, err)

// 	assert.True(t, len(descr.Modules) > 10)
// 	set := CreateModuleSet(createWorkspace(), descr.NormaliseModules("c:\\wf\\lp\\web\\App"), nil)
// 	assert.Equal(t, "controlPanel/ControlPanel", set.names()[0], "controlPanel/ControlPanel should be the first module")
// }
