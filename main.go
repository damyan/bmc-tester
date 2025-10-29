package main

import (
	"context"
	"log"
	"os"

	"github.com/damyan/bmc-tester/bmc"
	"github.com/spf13/cobra"
)

var redfishBMC *bmc.RedfishBMC
var config Config

type Config struct {
	Username          string
	Password          string
	Endpoint          string
	URISuffix         string
	EntityTag         string
	DisableEtagMatch  bool
	IfNoneMatchHeader string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "bmc-tester",
	Args: cobra.NoArgs,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	var err error
	if err = rootCmd.Execute(); err != nil {
		os.Exit(0)
	}

	if redfishBMC != nil {
		redfishBMC.Logout()
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&config.Username, "username", "u", "", "BMC username")
	rootCmd.PersistentFlags().StringVarP(&config.Password, "password", "p", "", "BMC password")
	rootCmd.PersistentFlags().StringVarP(&config.Endpoint, "endpoint", "e", "", "BMC endpoint")

	_ = rootCmd.MarkPersistentFlagRequired("username")
	_ = rootCmd.MarkPersistentFlagRequired("password")
	_ = rootCmd.MarkPersistentFlagRequired("endpoint")

	rootCmd.AddCommand(NewBootOnceCommand())
	rootCmd.AddCommand(NewPowerCommand())
}

func main() {
	Execute()
}

func NewBootOnceCommand() *cobra.Command {
	bootOnceCmd := &cobra.Command{
		Use:   "boot-once",
		Short: "Set/Disable boot once",
	}
	bootOnceCmd.PersistentFlags().StringVarP(&config.URISuffix, "uri-suffix", "s", "", "BMC URI suffix")
	bootOnceCmd.PersistentFlags().StringVarP(&config.EntityTag, "entity-tag", "t", "", "BMC entity tag (etag)")
	bootOnceCmd.PersistentFlags().StringVarP(&config.IfNoneMatchHeader, "if-none-match-header", "i", "", "Set if-none-match header")
	bootOnceCmd.PersistentFlags().BoolVarP(&config.DisableEtagMatch, "disable-etag-match", "d", false, "Disable etag match (no header)")
	bootOnceCmd.AddCommand(NewBootOncePXECommand())
	bootOnceCmd.AddCommand(NewBootOnceDisableCommand())
	bootOnceCmd.AddCommand(NewGetBootOnceCommand())
	return bootOnceCmd
}

func NewBootOncePXECommand() *cobra.Command {
	setBootOncePXECmd := &cobra.Command{
		Use:   "PXE",
		Short: "Set boot once to PXE",
		RunE:  runSetBootPXE,
	}
	return setBootOncePXECmd
}

func NewBootOnceDisableCommand() *cobra.Command {
	setBootOnceDisable := &cobra.Command{
		Use:   "disable",
		Short: "Disable boot once",
		RunE:  runSetBootDisable,
	}
	return setBootOnceDisable
}

func NewGetBootOnceCommand() *cobra.Command {
	setBootOnceDisable := &cobra.Command{
		Use:   "get",
		Short: "Get boot once",
		RunE:  runGetBootOnce,
	}
	return setBootOnceDisable
}

func runSetBootPXE(cmd *cobra.Command, args []string) error {
	log.Print("Setting boot once to PXE...")
	opt := initOptions()

	var err error
	redfishBMC, err = bmc.NewRedfishBMCClient(context.TODO(), *opt)

	if redfishBMC == nil {
		log.Fatalf("Could not create redfish bmc: %v", err)
	}

	return redfishBMC.RunSetBootOncePXE()
}

func runSetBootDisable(cmd *cobra.Command, args []string) error {
	log.Print("Disable boot once...")
	opt := initOptions()

	var err error
	redfishBMC, err = bmc.NewRedfishBMCClient(context.TODO(), *opt)

	if redfishBMC == nil {
		log.Fatalf("Could not create redfish bmc: %v", err)
	}

	return redfishBMC.RunSetBootOnceDisable()
}

func runGetBootOnce(cmd *cobra.Command, args []string) error {
	log.Print("Getting boot once...")
	opt := initOptions()

	var err error
	redfishBMC, err = bmc.NewRedfishBMCClient(context.TODO(), *opt)

	if redfishBMC == nil {
		log.Fatalf("Could not create redfish bmc: %v", err)
	}

	return redfishBMC.RunGetBootOnce()
}

func NewPowerCommand() *cobra.Command {
	powerCmd := &cobra.Command{
		Use:   "power",
		Short: "Power on/off",
	}
	powerCmd.PersistentFlags().StringVarP(&config.URISuffix, "uri-suffix", "s", "", "BMC URI suffix")
	powerCmd.PersistentFlags().StringVarP(&config.EntityTag, "entity-tag", "t", "", "BMC entity tag (etag)")
	powerCmd.PersistentFlags().StringVarP(&config.IfNoneMatchHeader, "if-none-match-header", "i", "", "Set if-none-match header")
	powerCmd.PersistentFlags().BoolVarP(&config.DisableEtagMatch, "disable-etag-match", "d", false, "Disable etag match (no header)")
	powerCmd.AddCommand(NewPowerOnCommand())
	powerCmd.AddCommand(NewPowerOffCommand())
	powerCmd.AddCommand(NewGetPowerCommand())
	return powerCmd
}

func NewPowerOnCommand() *cobra.Command {
	powerOnCmd := &cobra.Command{
		Use:   "on",
		Short: "Set power on",
		RunE:  runPowerOn,
	}
	return powerOnCmd
}

func NewPowerOffCommand() *cobra.Command {
	powerOnCmd := &cobra.Command{
		Use:   "off",
		Short: "Set power off",
		RunE:  runPowerOff,
	}
	return powerOnCmd
}

func NewGetPowerCommand() *cobra.Command {
	getPowerCmd := &cobra.Command{
		Use:   "get",
		Short: "Get power",
		RunE:  runGetPower,
	}
	return getPowerCmd
}

func runPowerOn(cmd *cobra.Command, args []string) error {
	log.Print("Setting power on...")
	opt := initOptions()

	var err error
	redfishBMC, err = bmc.NewRedfishBMCClient(context.TODO(), *opt)

	if redfishBMC == nil {
		log.Fatalf("Could not create redfish bmc: %v", err)
	}

	return redfishBMC.RunPowerOn()
}

func runPowerOff(cmd *cobra.Command, args []string) error {
	log.Print("Setting power off...")
	opt := initOptions()

	var err error
	redfishBMC, err = bmc.NewRedfishBMCClient(context.TODO(), *opt)

	if redfishBMC == nil {
		log.Fatalf("Could not create redfish bmc: %v", err)
	}

	return redfishBMC.RunPowerOff()
}

func runGetPower(cmd *cobra.Command, args []string) error {
	log.Print("Getting power...")
	opt := initOptions()

	var err error
	redfishBMC, err = bmc.NewRedfishBMCClient(context.TODO(), *opt)

	if redfishBMC == nil {
		log.Fatalf("Could not create redfish bmc: %v", err)
	}

	return redfishBMC.RunGetPower()
}

func initOptions() *bmc.Options {
	return &bmc.Options{
		Username:          config.Username,
		Password:          config.Password,
		Endpoint:          config.Endpoint,
		BasicAuth:         true,
		URISuffix:         config.URISuffix,
		EntityTag:         config.EntityTag,
		DisableEtagMatch:  config.DisableEtagMatch,
		IfNoneMatchHeader: config.IfNoneMatchHeader,
	}
}
