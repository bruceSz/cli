package route_test

import (
	"cf"
	. "cf/commands/route"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestListingRoutes(t *testing.T) {
	routes := []cf.Route{
		cf.Route{
			Host:     "hostname-1",
			Domain:   cf.Domain{Name: "example.com"},
			AppNames: []string{"dora", "dora2"},
		},
		cf.Route{
			Host:     "hostname-2",
			Domain:   cf.Domain{Name: "cfapps.com"},
			AppNames: []string{"my-app", "my-app2"},
		},
		cf.Route{
			Host:     "hostname-3",
			Domain:   cf.Domain{Name: "another-example.com"},
			AppNames: []string{"july", "june"},
		},
	}

	routeRepo := &testapi.FakeRouteRepository{Routes: routes}

	ui := callListRoutes(t, []string{}, &testreq.FakeReqFactory{}, routeRepo)

	assert.Contains(t, ui.Outputs[0], "Getting routes")
	assert.Contains(t, ui.Outputs[0], "my-user")

	assert.Contains(t, ui.Outputs[1], "host")
	assert.Contains(t, ui.Outputs[1], "domain")
	assert.Contains(t, ui.Outputs[1], "apps")

	assert.Contains(t, ui.Outputs[2], "hostname-1")
	assert.Contains(t, ui.Outputs[2], "example.com")
	assert.Contains(t, ui.Outputs[2], "dora, dora2")

	assert.Contains(t, ui.Outputs[3], "hostname-2")
	assert.Contains(t, ui.Outputs[3], "cfapps.com")
	assert.Contains(t, ui.Outputs[3], "my-app, my-app2")

	assert.Contains(t, ui.Outputs[4], "hostname-3")
	assert.Contains(t, ui.Outputs[4], "another-example.com")
	assert.Contains(t, ui.Outputs[4], "july, june")
}

func TestListingRoutesWhenNoneExist(t *testing.T) {
	routes := []cf.Route{}
	routeRepo := &testapi.FakeRouteRepository{Routes: routes}

	ui := callListRoutes(t, []string{}, &testreq.FakeReqFactory{}, routeRepo)

	assert.Contains(t, ui.Outputs[0], "Getting routes")
	assert.Contains(t, ui.Outputs[0], "my-user")
	assert.Contains(t, ui.Outputs[1], "No routes found")
}

func TestListingRoutesWhenFindFails(t *testing.T) {
	routeRepo := &testapi.FakeRouteRepository{ListErr: true}

	ui := callListRoutes(t, []string{}, &testreq.FakeReqFactory{}, routeRepo)

	assert.Contains(t, ui.Outputs[0], "Getting routes")
	assert.Contains(t, ui.Outputs[1], "FAILED")
}

func callListRoutes(t *testing.T, args []string, reqFactory *testreq.FakeReqFactory, routeRepo *testapi.FakeRouteRepository) (ui *testterm.FakeUI) {

	ui = &testterm.FakeUI{}

	ctxt := testcmd.NewContext("list-routes", args)

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "my-user",
	})
	assert.NoError(t, err)

	config := &configuration.Configuration{
		Space:        cf.Space{Name: "my-space"},
		Organization: cf.Organization{Name: "my-org"},
		AccessToken:  token,
	}

	cmd := NewListRoutes(ui, config, routeRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)

	return
}
