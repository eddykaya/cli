package cmd

import (
	"errors"
	"fmt"
	"strings"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/spf13/cobra"
)

type sksNodepoolUpdateCmd struct {
	cliCommandSettings `cli-cmd:"-"`

	_ bool `cli-cmd:"update"`

	Cluster  string `cli-arg:"#" cli-usage:"CLUSTER-NAME|ID"`
	Nodepool string `cli-arg:"#" cli-usage:"NODEPOOL-NAME|ID"`

	AntiAffinityGroups []string `cli-flag:"anti-affinity-group" cli-usage:"Nodepool Anti-Affinity Group NAME|ID (can be specified multiple times)"`
	DeployTarget       string   `cli-usage:"Nodepool Deploy Target NAME|ID"`
	Description        string   `cli-usage:"Nodepool description"`
	DiskSize           int64    `cli-usage:"Nodepool Compute instances disk size"`
	InstancePrefix     string   `cli-usage:"string to prefix Nodepool member names with"`
	InstanceType       string   `cli-usage:"Nodepool Compute instances type"`
	Labels             []string `cli-flag:"label" cli-usage:"Nodepool label (format: KEY=VALUE, can be repeated multiple times)"`
	Name               string   `cli-usage:"Nodepool name"`
	PrivateNetworks    []string `cli-flag:"private-network" cli-usage:"Nodepool Private Network NAME|ID (can be specified multiple times)"`
	SecurityGroups     []string `cli-flag:"security-group" cli-usage:"Nodepool Security Group NAME|ID (can be specified multiple times)"`
	Taints             []string `cli-flag:"taint" cli-usage:"Kubernetes taint to apply to Nodepool Nodes (format: KEY=VALUE:EFFECT, can be specified multiple times)"`
	Zone               string   `cli-short:"z" cli-usage:"SKS cluster zone"`
}

func (c *sksNodepoolUpdateCmd) cmdAliases() []string { return nil }

func (c *sksNodepoolUpdateCmd) cmdShort() string { return "Update an SKS cluster Nodepool" }

func (c *sksNodepoolUpdateCmd) cmdLong() string {
	return fmt.Sprintf(`This command updates an SKS Nodepool.

Supported output template annotations: %s`,
		strings.Join(outputterTemplateAnnotations(&sksNodepoolShowOutput{}), ", "),
	)
}

func (c *sksNodepoolUpdateCmd) cmdPreRun(cmd *cobra.Command, args []string) error {
	cmdSetZoneFlagFromDefault(cmd)
	return cliCommandDefaultPreRun(c, cmd, args)
}

func (c *sksNodepoolUpdateCmd) cmdRun(cmd *cobra.Command, _ []string) error {
	var (
		nodepool *egoscale.SKSNodepool
		updated  bool
	)

	ctx := exoapi.WithEndpoint(gContext, exoapi.NewReqEndpoint(gCurrentAccount.Environment, c.Zone))

	cluster, err := cs.FindSKSCluster(ctx, c.Zone, c.Cluster)
	if err != nil {
		return err
	}

	for _, n := range cluster.Nodepools {
		if *n.ID == c.Nodepool || *n.Name == c.Nodepool {
			nodepool = n
			break
		}
	}
	if nodepool == nil {
		return errors.New("Nodepool not found") // nolint:golint
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.AntiAffinityGroups)) {
		nodepoolAntiAffinityGroupIDs := make([]string, len(c.AntiAffinityGroups))
		for i, v := range c.AntiAffinityGroups {
			antiAffinityGroup, err := cs.FindAntiAffinityGroup(ctx, c.Zone, v)
			if err != nil {
				return fmt.Errorf("error retrieving Anti-Affinity Group: %w", err)
			}
			nodepoolAntiAffinityGroupIDs[i] = *antiAffinityGroup.ID
		}
		nodepool.AntiAffinityGroupIDs = &nodepoolAntiAffinityGroupIDs
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.DeployTarget)) {
		deployTarget, err := cs.FindDeployTarget(ctx, c.Zone, c.DeployTarget)
		if err != nil {
			return fmt.Errorf("error retrieving Deploy Target: %w", err)
		}
		nodepool.DeployTargetID = deployTarget.ID
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.Description)) {
		nodepool.Description = &c.Description
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.DiskSize)) {
		nodepool.DiskSize = &c.DiskSize
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.InstancePrefix)) {
		nodepool.InstancePrefix = &c.InstancePrefix
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.InstanceType)) {
		nodepoolInstanceType, err := cs.FindInstanceType(ctx, c.Zone, c.InstanceType)
		if err != nil {
			return fmt.Errorf("error retrieving instance type: %w", err)
		}
		nodepool.InstanceTypeID = nodepoolInstanceType.ID
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.Labels)) {
		labels := make(map[string]string)
		if len(c.Labels) > 0 {
			labels, err = sliceToMap(c.Labels)
			if err != nil {
				return fmt.Errorf("label: %w", err)
			}
		}
		nodepool.Labels = &labels
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.Name)) {
		nodepool.Name = &c.Name
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.PrivateNetworks)) {
		nodepoolPrivateNetworkIDs := make([]string, len(c.PrivateNetworks))
		for i, v := range c.PrivateNetworks {
			privateNetwork, err := cs.FindPrivateNetwork(ctx, c.Zone, v)
			if err != nil {
				return fmt.Errorf("error retrieving Private Network: %w", err)
			}
			nodepoolPrivateNetworkIDs[i] = *privateNetwork.ID
		}
		nodepool.PrivateNetworkIDs = &nodepoolPrivateNetworkIDs
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.SecurityGroups)) {
		nodepoolSecurityGroupIDs := make([]string, len(c.SecurityGroups))
		for i, v := range c.SecurityGroups {
			securityGroup, err := cs.FindSecurityGroup(ctx, c.Zone, v)
			if err != nil {
				return fmt.Errorf("error retrieving Security Group: %w", err)
			}
			nodepoolSecurityGroupIDs[i] = *securityGroup.ID
		}
		nodepool.SecurityGroupIDs = &nodepoolSecurityGroupIDs
		updated = true
	}

	if cmd.Flags().Changed(mustCLICommandFlagName(c, &c.Taints)) {
		taints := make(map[string]*egoscale.SKSNodepoolTaint)
		for _, t := range c.Taints {
			key, taint, err := parseSKSNodepoolTaint(t)
			if err != nil {
				return fmt.Errorf("invalid taint value %q: %w", t, err)
			}
			taints[key] = taint
		}
		nodepool.Taints = &taints
		updated = true
	}

	if updated {
		decorateAsyncOperation(fmt.Sprintf("Updating Nodepool %q...", c.Nodepool), func() {
			if err = cs.UpdateSKSNodepool(ctx, c.Zone, cluster, nodepool); err != nil {
				return
			}
		})
		if err != nil {
			return err
		}
	}

	if !gQuiet {
		return (&sksNodepoolShowCmd{
			cliCommandSettings: c.cliCommandSettings,
			Cluster:            *cluster.ID,
			Nodepool:           *nodepool.ID,
			Zone:               c.Zone,
		}).cmdRun(nil, nil)
	}

	return nil
}

func init() {
	cobra.CheckErr(registerCLICommand(sksNodepoolCmd, &sksNodepoolUpdateCmd{
		cliCommandSettings: defaultCLICmdSettings(),
	}))
}
