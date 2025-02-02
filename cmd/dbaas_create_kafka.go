package cmd

import (
	"fmt"
	"net/http"

	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/spf13/cobra"
)

func (c *dbaasServiceCreateCmd) createKafka(_ *cobra.Command, _ []string) error {
	var err error

	ctx := exoapi.WithEndpoint(gContext, exoapi.NewReqEndpoint(gCurrentAccount.Environment, c.Zone))

	databaseService := oapi.CreateDbaasServiceKafkaJSONRequestBody{
		KafkaConnectEnabled:   &c.KafkaEnableKafkaConnect,
		KafkaRestEnabled:      &c.KafkaEnableKafkaREST,
		Plan:                  c.Plan,
		SchemaRegistryEnabled: &c.KafkaEnableSchemaRegistry,
		TerminationProtection: &c.TerminationProtection,
		Version:               nonEmptyStringPtr(c.KafkaVersion),
	}

	settingsSchema, err := cs.GetDbaasSettingsKafkaWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve Database Service settings: %w", err)
	}
	if settingsSchema.StatusCode() != http.StatusOK {
		return fmt.Errorf("API request error: unexpected status %s", settingsSchema.Status())
	}

	if c.KafkaEnableCertAuth || c.KafkaEnableSASLAuth {
		databaseService.AuthenticationMethods = &struct {
			Certificate *bool `json:"certificate,omitempty"`
			Sasl        *bool `json:"sasl,omitempty"`
		}{
			Certificate: &c.KafkaEnableCertAuth,
			Sasl:        &c.KafkaEnableSASLAuth,
		}
	}

	if len(c.KafkaIPFilter) > 0 {
		databaseService.IpFilter = &c.KafkaIPFilter
	}

	if c.MaintenanceDOW != "" && c.MaintenanceTime != "" {
		databaseService.Maintenance = &struct {
			Dow  oapi.CreateDbaasServiceKafkaJSONBodyMaintenanceDow `json:"dow"`
			Time string                                             `json:"time"`
		}{
			Dow:  oapi.CreateDbaasServiceKafkaJSONBodyMaintenanceDow(c.MaintenanceDOW),
			Time: c.MaintenanceTime,
		}
	}

	if c.KafkaConnectSettings != "" {
		settings, err := validateDatabaseServiceSettings(
			c.KafkaConnectSettings,
			settingsSchema.JSON200.Settings.KafkaConnect,
		)
		if err != nil {
			return fmt.Errorf("invalid settings: %w", err)
		}
		databaseService.KafkaConnectSettings = &settings
	}

	if c.KafkaRESTSettings != "" {
		settings, err := validateDatabaseServiceSettings(
			c.KafkaRESTSettings,
			settingsSchema.JSON200.Settings.KafkaRest,
		)
		if err != nil {
			return fmt.Errorf("invalid settings: %w", err)
		}
		databaseService.KafkaRestSettings = &settings
	}

	if c.KafkaSettings != "" {
		settings, err := validateDatabaseServiceSettings(
			c.KafkaSettings,
			settingsSchema.JSON200.Settings.Kafka,
		)
		if err != nil {
			return fmt.Errorf("invalid settings: %w", err)
		}
		databaseService.KafkaSettings = &settings
	}

	if c.KafkaSchemaRegistrySettings != "" {
		settings, err := validateDatabaseServiceSettings(
			c.KafkaSchemaRegistrySettings,
			settingsSchema.JSON200.Settings.SchemaRegistry,
		)
		if err != nil {
			return fmt.Errorf("invalid settings: %w", err)
		}
		databaseService.SchemaRegistrySettings = &settings
	}

	var res *oapi.CreateDbaasServiceKafkaResponse
	decorateAsyncOperation(fmt.Sprintf("Creating Database Service %q...", c.Name), func() {
		res, err = cs.CreateDbaasServiceKafkaWithResponse(ctx, oapi.DbaasServiceName(c.Name), databaseService)
	})
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("API request error: unexpected status %s", res.Status())
	}

	if !gQuiet {
		return c.outputFunc((&dbaasServiceShowCmd{
			Name: c.Name,
			Zone: c.Zone,
		}).showDatabaseServiceKafka(ctx))
	}

	return nil
}
