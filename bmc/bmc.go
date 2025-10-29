package bmc

import (
	"context"
	"fmt"
	"log"

	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
	//"github.com/stmcginnis/gofish/common"
	//"github.com/stmcginnis/gofish/redfish"
)

var pxeBootWithSettingUEFIBootMode = redfish.Boot{
	BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	BootSourceOverrideMode:    redfish.UEFIBootSourceOverrideMode,
	BootSourceOverrideTarget:  redfish.PxeBootSourceOverrideTarget,
}
var pxeBootWithoutSettingUEFIBootMode = redfish.Boot{
	BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	BootSourceOverrideTarget:  redfish.PxeBootSourceOverrideTarget,
}

var disableBootWithSettingUEFIBootMode = redfish.Boot{
	BootSourceOverrideEnabled: redfish.DisabledBootSourceOverrideEnabled,
	BootSourceOverrideMode:    redfish.UEFIBootSourceOverrideMode,
}
var disableBootWithoutSettingUEFIBootMode = redfish.Boot{
	BootSourceOverrideEnabled: redfish.DisabledBootSourceOverrideEnabled,
	BootSourceOverrideTarget:  redfish.PxeBootSourceOverrideTarget,
}

type Options struct {
	Endpoint         string
	Username         string
	Password         string
	BasicAuth        bool
	URISuffix        string
	EntityTag        string
	DisableEtagMatch bool
}

type RedfishBMC struct {
	client  *gofish.APIClient
	Options Options
}

func NewRedfishBMCClient(ctx context.Context, options Options) (*RedfishBMC, error) {
	clientConfig := gofish.ClientConfig{
		Endpoint:  options.Endpoint,
		Username:  options.Username,
		Password:  options.Password,
		Insecure:  true,
		BasicAuth: options.BasicAuth,
	}
	client, err := gofish.ConnectContext(ctx, clientConfig)
	if err != nil {
		return nil, err
	}
	bmc := &RedfishBMC{client: client, Options: options}

	return bmc, nil
}

func (r *RedfishBMC) RunSetBootOncePXE() error {
	systems, err := r.getSystems()
	if err != nil {
		return err
	}

	for _, system := range systems {
		var setBoot redfish.Boot
		if system.Boot.BootSourceOverrideMode != redfish.UEFIBootSourceOverrideMode {
			setBoot = pxeBootWithSettingUEFIBootMode
		} else {
			setBoot = pxeBootWithoutSettingUEFIBootMode
		}
		if err := system.SetBoot(setBoot); err != nil {
			return fmt.Errorf("failed to set the boot order: %w", err)
		}
	}

	return nil
}

func (r *RedfishBMC) RunSetBootOnceDisable() error {
	systems, err := r.getSystems()
	if err != nil {
		return err
	}

	for _, system := range systems {
		var setBoot redfish.Boot
		if system.Boot.BootSourceOverrideMode != redfish.UEFIBootSourceOverrideMode {
			setBoot = disableBootWithSettingUEFIBootMode
		} else {
			setBoot = disableBootWithoutSettingUEFIBootMode
		}
		if err := system.SetBoot(setBoot); err != nil {
			return fmt.Errorf("failed to set the boot order: %w", err)
		}
	}

	return nil
}

func (r *RedfishBMC) RunGetBootOnce() error {
	systems, err := r.getSystems()
	if err != nil {
		return err
	}
	for _, system := range systems {
		var getBoot redfish.Boot
		getBoot = system.Boot

		fmt.Printf("Enabled: %s, Target: %s\n", getBoot.BootSourceOverrideEnabled, getBoot.BootSourceOverrideTarget)
	}
	return nil
}

func (r *RedfishBMC) RunPowerOn() error {
	systems, err := r.getSystems()
	if err != nil {
		return err
	}

	for _, system := range systems {
		powerState := system.PowerState
		if powerState != redfish.OnPowerState {
			if err := system.Reset(redfish.OnResetType); err != nil {
				return fmt.Errorf("failed to reset system to power on state: %w", err)
			}
		}
	}

	return nil
}

func (r *RedfishBMC) RunPowerOff() error {
	systems, err := r.getSystems()
	if err != nil {
		return err
	}

	for _, system := range systems {
		powerState := system.PowerState
		if powerState != redfish.OffPowerState {
			if err := system.Reset(redfish.ForceOffResetType); err != nil {
				return fmt.Errorf("failed to reset system to power off state: %w", err)
			}
		}
	}

	return nil
}

func (r *RedfishBMC) RunGetPower() error {
	systems, err := r.getSystems()
	if err != nil {
		return err
	}

	for _, system := range systems {
		powerState := system.PowerState
		fmt.Printf("Power: %s\n", powerState)
	}

	return nil
}

func (r *RedfishBMC) Logout() {
	if r.client != nil {
		r.client.Logout()
	}
}

func (r *RedfishBMC) getSystems() ([]*redfish.ComputerSystem, error) {
	service := r.client.GetService()
	systems, err := service.Systems()

	for _, system := range systems {
		if r.Options.EntityTag != "" {
			system.SetETag(r.Options.EntityTag)
		}

		system.DisableEtagMatch(r.Options.DisableEtagMatch)

		if r.Options.URISuffix != "" {
			system.ODataID = system.ODataID + r.Options.URISuffix
		}
		log.Printf("System URI: %s", system.ODataID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get systems: %w", err)
	}
	return systems, nil
}

func (r *RedfishBMC) getSystemFromUri() (*redfish.ComputerSystem, error) {

	var systemURI string

	if r.Options.URISuffix != "" {
		systemURI = r.Options.URISuffix
	}

	system, err := common.GetObject[redfish.ComputerSystem](r.client, systemURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get system: %w", err)
	}

	if system.UUID != "" {
		return system, nil
	}

	return nil, fmt.Errorf("no system found for %w", err)
}
