package servicebroker

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/codegangsta/cli"
)

type CreateServiceBroker struct {
	ui                terminal.UI
	config            *configuration.Configuration
	serviceBrokerRepo api.ServiceBrokerRepository
}

func NewCreateServiceBroker(ui terminal.UI, config *configuration.Configuration, serviceBrokerRepo api.ServiceBrokerRepository) (cmd CreateServiceBroker) {
	cmd.ui = ui
	cmd.config = config
	cmd.serviceBrokerRepo = serviceBrokerRepo
	return
}

func (cmd CreateServiceBroker) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {

	if len(c.Args()) != 4 {
		err = errors.New("Incorrect usage")
		cmd.ui.FailWithUsage(c, "create-service-broker")
		return
	}

	reqs = append(reqs, reqFactory.NewLoginRequirement())

	return
}

func (cmd CreateServiceBroker) Run(c *cli.Context) {
	serviceBroker := cf.ServiceBroker{
		Name:     c.Args()[0],
		Username: c.Args()[1],
		Password: c.Args()[2],
		Url:      c.Args()[3],
	}
	cmd.ui.Say("Creating service broker %s as %s...",
		terminal.EntityNameColor(serviceBroker.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiResponse := cmd.serviceBrokerRepo.Create(serviceBroker)
	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
		return
	}

	cmd.ui.Ok()
}
