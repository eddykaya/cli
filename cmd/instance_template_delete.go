package cmd

import (
	"fmt"

	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/spf13/cobra"
)

type instanceTemplateDeleteCmd struct {
	cliCommandSettings `cli-cmd:"-"`

	_ bool `cli-cmd:"delete"`

	TemplateID string `cli-arg:"#" cli-usage:"TEMPLATE-ID"`

	Force bool   `cli-short:"f" cli-usage:"don't prompt for confirmation"`
	Zone  string `cli-short:"z" cli-usage:"template zone"`
}

func (c *instanceTemplateDeleteCmd) cmdAliases() []string { return gRemoveAlias }

func (c *instanceTemplateDeleteCmd) cmdShort() string {
	return "Delete a Compute instance template"
}

func (c *instanceTemplateDeleteCmd) cmdLong() string { return "" }

func (c *instanceTemplateDeleteCmd) cmdPreRun(cmd *cobra.Command, args []string) error {
	cmdSetZoneFlagFromDefault(cmd)
	return cliCommandDefaultPreRun(c, cmd, args)
}

func (c *instanceTemplateDeleteCmd) cmdRun(_ *cobra.Command, _ []string) error {
	ctx := exoapi.WithEndpoint(gContext, exoapi.NewReqEndpoint(gCurrentAccount.Environment, c.Zone))

	template, err := cs.GetTemplate(ctx, c.Zone, c.TemplateID)
	if err != nil {
		return err
	}

	if !c.Force {
		if !askQuestion(fmt.Sprintf(
			"Are you sure you want to delete template %s (%q)?",
			c.TemplateID,
			*template.Name,
		)) {
			return nil
		}
	}

	decorateAsyncOperation(fmt.Sprintf("Deleting template %s...", c.TemplateID), func() {
		err = cs.DeleteTemplate(ctx, c.Zone, template)
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	cobra.CheckErr(registerCLICommand(instanceTemplateCmd, &instanceTemplateDeleteCmd{
		cliCommandSettings: defaultCLICmdSettings(),
	}))
}
