package organization

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/codegangsta/cli"
)

type DeleteOrg struct {
	ui         terminal.UI
	config     *configuration.Configuration
	orgRepo    api.OrganizationRepository
	orgReq     requirements.OrganizationRequirement
	configRepo configuration.ConfigurationRepository
}

func NewDeleteOrg(ui terminal.UI, config *configuration.Configuration, sR api.OrganizationRepository, cR configuration.ConfigurationRepository) (cmd *DeleteOrg) {
	cmd = new(DeleteOrg)
	cmd.ui = ui
	cmd.config = config
	cmd.orgRepo = sR
	cmd.configRepo = cR
	return
}

func (cmd *DeleteOrg) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "delete-org")
		return

	}
	return
}

func (cmd *DeleteOrg) Run(c *cli.Context) {
	orgName := c.Args()[0]

	force := c.Bool("f")

	if !force {
		response := cmd.ui.Confirm(
			"Really delete org %s and everything associated with it?%s",
			terminal.EntityNameColor(orgName),
			terminal.PromptColor(">"),
		)

		if !response {
			return
		}
	}

	cmd.ui.Say("Deleting org %s as %s...",
		terminal.EntityNameColor(orgName),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	org, apiResponse := cmd.orgRepo.FindByName(orgName)

	if apiResponse.IsError() {
		cmd.ui.Failed(apiResponse.Message)
		return
	}

	if apiResponse.IsNotFound() {
		cmd.ui.Ok()
		cmd.ui.Warn("Org %s does not exist.", orgName)
		return
	}

	apiResponse = cmd.orgRepo.Delete(org)
	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
		return
	}
	config, err := cmd.configRepo.Get()
	if err != nil {
		cmd.ui.Failed("Couldn't reset your target. You should logout and log in again.")
		return
	}

	if org.Guid == config.Organization.Guid {
		config.Organization = cf.Organization{}
		config.Space = cf.Space{}
		cmd.configRepo.Save()
	}

	cmd.ui.Ok()
	return
}
