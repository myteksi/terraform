package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/hashicorp/terraform/command"
	pluginDiscovery "github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/svchost"
	"github.com/hashicorp/terraform/svchost/auth"
	"github.com/hashicorp/terraform/svchost/disco"
	"github.com/mitchellh/cli"
)

// runningInAutomationEnvName gives the name of an environment variable that
// can be set to any non-empty value in order to suppress certain messages
// that assume that Terraform is being run from a command prompt.
const runningInAutomationEnvName = "TF_IN_AUTOMATION"

// Commands is the mapping of all the available Terraform commands.
var Commands map[string]cli.CommandFactory
var PlumbingCommands map[string]struct{}

// Ui is the cli.Ui used for communicating to the outside world.
var Ui cli.Ui

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
)

func initCommands(config *Config) {
	var inAutomation bool
	if v := os.Getenv(runningInAutomationEnvName); v != "" {
		inAutomation = true
	}

	credsSrc := credentialsSource(config)
	services := disco.NewDisco()
	services.SetCredentialsSource(credsSrc)
	for userHost, hostConfig := range config.Hosts {
		host, err := svchost.ForComparison(userHost)
		if err != nil {
			// We expect the config was already validated by the time we get
			// here, so we'll just ignore invalid hostnames.
			continue
		}
		services.ForceHostServices(host, hostConfig.Services)
	}

	dataDir := os.Getenv("TF_DATA_DIR")

	meta := command.Meta{
		Color:            true,
		GlobalPluginDirs: globalPluginDirs(),
		PluginOverrides:  &PluginOverrides,
		Ui:               Ui,

		Services:    services,
		Credentials: credsSrc,

		RunningInAutomation: inAutomation,
		PluginCacheDir:      config.PluginCacheDir,
		OverrideDataDir:     dataDir,

		ShutdownCh: makeShutdownCh(),
	}

	// The command list is included in the terraform -help
	// output, which is in turn included in the docs at
	// website/source/docs/commands/index.html.markdown; if you
	// add, remove or reclassify commands then consider updating
	// that to match.

	PlumbingCommands = map[string]struct{}{
		"state":        struct{}{}, // includes all subcommands
		"debug":        struct{}{}, // includes all subcommands
		"force-unlock": struct{}{},
	}

	Commands = map[string]cli.CommandFactory{
		"plan": func() (cli.Command, error) {
			return &command.PlanCommand{
				Meta: meta,
			}, nil
		},


	}
}

// makeShutdownCh creates an interrupt listener and returns a channel.
// A message will be sent on the channel for every interrupt received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, ignoreSignals...)
	signal.Notify(signalCh, forwardSignals...)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}

func credentialsSource(config *Config) auth.CredentialsSource {
	creds := auth.NoCredentials
	if len(config.Credentials) > 0 {
		staticTable := map[svchost.Hostname]map[string]interface{}{}
		for userHost, creds := range config.Credentials {
			host, err := svchost.ForComparison(userHost)
			if err != nil {
				// We expect the config was already validated by the time we get
				// here, so we'll just ignore invalid hostnames.
				continue
			}
			staticTable[host] = creds
		}
		creds = auth.StaticCredentialsSource(staticTable)
	}

	for helperType, helperConfig := range config.CredentialsHelpers {
		log.Printf("[DEBUG] Searching for credentials helper named %q", helperType)
		available := pluginDiscovery.FindPlugins("credentials", globalPluginDirs())
		available = available.WithName(helperType)
		if available.Count() == 0 {
			log.Printf("[ERROR] Unable to find credentials helper %q; ignoring", helperType)
			break
		}

		selected := available.Newest()

		helperSource := auth.HelperProgramCredentialsSource(selected.Path, helperConfig.Args...)
		creds = auth.Credentials{
			creds,
			auth.CachingCredentialsSource(helperSource), // cached because external operation may be slow/expensive
		}

		// There should only be zero or one "credentials_helper" blocks. We
		// assume that the config was validated earlier and so we don't check
		// for extras here.
		break
	}

	return creds
}
