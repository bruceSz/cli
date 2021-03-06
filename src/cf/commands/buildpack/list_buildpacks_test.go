package buildpack_test

import (
	"cf"
	. "cf/commands/buildpack"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestListBuildpacksRequirements(t *testing.T) {
	buildpackRepo := &testapi.FakeBuildpackRepository{}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	callListBuildpacks(reqFactory, buildpackRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: false}
	callListBuildpacks(reqFactory, buildpackRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestListBuildpacks(t *testing.T) {
	position5 := 5
	position10 := 10
	position15 := 15
	buildpacks := []cf.Buildpack{
		{Name: "Buildpack-1", Position: &position5},
		{Name: "Buildpack-2", Position: &position10},
		{Name: "Buildpack-3", Position: &position15},
	}

	buildpackRepo := &testapi.FakeBuildpackRepository{
		Buildpacks: buildpacks,
	}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

	ui := callListBuildpacks(reqFactory, buildpackRepo)

	assert.Contains(t, ui.Outputs[0], "Getting buildpacks")

	assert.Contains(t, ui.Outputs[1], "buildpack")
	assert.Contains(t, ui.Outputs[1], "position")

	assert.Contains(t, ui.Outputs[2], "Buildpack-1")
	assert.Contains(t, ui.Outputs[2], "5")

	assert.Contains(t, ui.Outputs[3], "Buildpack-2")
	assert.Contains(t, ui.Outputs[3], "10")

	assert.Contains(t, ui.Outputs[4], "Buildpack-3")
	assert.Contains(t, ui.Outputs[4], "15")
}

func TestListingBuildpacksWhenNoneExist(t *testing.T) {
	buildpacks := []cf.Buildpack{}
	buildpackRepo := &testapi.FakeBuildpackRepository{
		Buildpacks: buildpacks,
	}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

	ui := callListBuildpacks(reqFactory, buildpackRepo)

	assert.Contains(t, ui.Outputs[0], "Getting buildpacks")
	assert.Contains(t, ui.Outputs[1], "No buildpacks found")
}

func callListBuildpacks(reqFactory *testreq.FakeReqFactory, buildpackRepo *testapi.FakeBuildpackRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = &testterm.FakeUI{}
	ctxt := testcmd.NewContext("buildpacks", []string{})
	cmd := NewListBuildpacks(fakeUI, buildpackRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
