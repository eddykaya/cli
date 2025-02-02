package cmd

import (
	"errors"
	"fmt"

	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/spf13/cobra"
)

type sksNodepoolDeleteCmd struct {
	cliCommandSettings `cli-cmd:"-"`

	_ bool `cli-cmd:"delete"`

	Cluster  string `cli-arg:"#" cli-usage:"CLUSTER-NAME|ID"`
	Nodepool string `cli-arg:"#" cli-usage:"NODEPOOL-NAME|ID"`

	Force bool   `cli-short:"f" cli-usage:"don't prompt for confirmation"`
	Zone  string `cli-short:"z" cli-usage:"SKS cluster zone"`
}

func (c *sksNodepoolDeleteCmd) cmdAliases() []string { return gRemoveAlias }

func (c *sksNodepoolDeleteCmd) cmdShort() string { return "Delete an SKS cluster Nodepool" }

func (c *sksNodepoolDeleteCmd) cmdLong() string { return "" }

func (c *sksNodepoolDeleteCmd) cmdPreRun(cmd *cobra.Command, args []string) error {
	cmdSetZoneFlagFromDefault(cmd)
	return cliCommandDefaultPreRun(c, cmd, args)
}

func (c *sksNodepoolDeleteCmd) cmdRun(_ *cobra.Command, _ []string) error {
	if !c.Force {
		if !askQuestion(fmt.Sprintf("Are you sure you want to delete Nodepool %q?", c.Nodepool)) {
			return nil
		}
	}

	ctx := exoapi.WithEndpoint(gContext, exoapi.NewReqEndpoint(gCurrentAccount.Environment, c.Zone))

	cluster, err := cs.FindSKSCluster(ctx, c.Zone, c.Cluster)
	if err != nil {
		return err
	}

	for _, nodepool := range cluster.Nodepools {
		if *nodepool.ID == c.Nodepool || *nodepool.Name == c.Nodepool {
			nodepool := nodepool
			decorateAsyncOperation(fmt.Sprintf("Deleting Nodepool %q...", *nodepool.Name), func() {
				err = cs.DeleteSKSNodepool(ctx, c.Zone, cluster, nodepool)
			})
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("Nodepool not found") // nolint:golint
}

func init() {
	cobra.CheckErr(registerCLICommand(sksNodepoolCmd, &sksNodepoolDeleteCmd{
		cliCommandSettings: defaultCLICmdSettings(),
	}))
}
