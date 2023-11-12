package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/catalog"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/gophercloud/openstack/placement/v1/resourceproviders"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const (
	//nolint:varnamelen
	OK          = "OK"
	Fail        = "Fail"
	Testing     = "Testing"
	Unsupported = "Unsupported"
)

var stateColors = map[string]termenv.ANSIColor{ //nolint:gochecknoglobals
	OK:          termenv.ANSIBrightGreen,
	Fail:        termenv.ANSIBrightRed,
	Testing:     termenv.ANSIBrightCyan,
	Unsupported: termenv.ANSIYellow,
}

// NewSmokeTestsCommand creates the root command.
func NewSmokeTestsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smoke-tests [CLOUD]",
		Short: "Run smoke tests",
		Long:  "Run smoke tests",
		Run:   smokeTestsCommand,
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllCloudNames)(cmd, args, toComplete)
			}
		}(),
	}

	return cmd
}

type regionName string

type serviceType string

type serviceName string

type serviceStateMsg struct {
	region  regionName
	service serviceName
	state   string
	err     string
}

type serviceStates map[regionName]map[serviceName]string

func serviceClient(regName regionName, svcType string) (*gophercloud.ServiceClient, context.CancelFunc) {
	opts := &clientconfig.ClientOpts{
		RegionName: string(regName),
		YAMLOpts:   config.YAMLOpts{Directory: viper.GetString("os-config-dir")},
	}

	client, err := clientconfig.NewServiceClient(svcType, opts)
	if err != nil {
		panic(err)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)

	client.Context = ctx

	return client, cancelFn
}

//nolint:funlen,cyclop
func smokeTest(regName regionName, svcName serviceName) tea.Msg {
	var (
		err       error
		cancel    context.CancelFunc
		svcClient *gophercloud.ServiceClient
	)

	switch string(svcName) {
	case "nova":
		svcClient, cancel = serviceClient(regName, "compute")

		_, err = servers.List(svcClient, servers.ListOpts{}).AllPages()
	case "neutron":
		svcClient, cancel = serviceClient(regName, "network")

		_, err = networks.List(svcClient, networks.ListOpts{}).AllPages()
	case "cinder":
		svcClient, cancel = serviceClient(regName, "volume")

		_, err = volumes.List(svcClient, volumes.ListOpts{}).AllPages()
	case "glance":
		svcClient, cancel = serviceClient(regName, "image")

		_, err = images.List(svcClient, images.ListOpts{}).AllPages()
		if err != nil && err.Error() == "json: cannot unmarshal string into Go struct field .images of type bool" {
			// some clouds sometimes return os_hidden as a string instead of a bool
			err = nil
		}
	case "heat":
		svcClient, cancel = serviceClient(regName, "orchestration")

		_, err = stacks.List(svcClient, stacks.ListOpts{}).AllPages()
	case "keystone":
		svcClient, cancel = serviceClient(regName, "identity")

		_, err = regions.List(svcClient, regions.ListOpts{}).AllPages()
	case "octavia":
		svcClient, cancel = serviceClient(regName, "load-balancer")

		_, err = loadbalancers.List(svcClient, loadbalancers.ListOpts{}).AllPages()
	case "swift":
		svcClient, cancel = serviceClient(regName, "object-store")

		_, err = containers.List(svcClient, containers.ListOpts{}).AllPages()
	case "placement":
		svcClient, cancel = serviceClient(regName, "placement")

		_, err = resourceproviders.List(svcClient, resourceproviders.ListOpts{}).AllPages()
	case "barbican":
		svcClient, cancel = serviceClient(regName, "key-manager")

		_, err = secrets.List(svcClient, secrets.ListOpts{}).AllPages()
	default:
		return serviceStateMsg{
			region:  regName,
			service: svcName,
			state:   Unsupported,
		}
	}

	defer cancel()

	if err == nil {
		return serviceStateMsg{
			region:  regName,
			service: svcName,
			state:   OK,
		}
	}

	return serviceStateMsg{
		region:  regName,
		service: svcName,
		state:   Fail,
		err:     err.Error(),
	}
}

type errorMsg struct {
	region  string
	service string
	message string
}

func (e errorMsg) String() string {
	if e.service != "" && e.region != "" {
		return fmt.Sprintf("%v(%v): %v", e.service, e.region, e.message)
	}

	return e.message
}

type model struct {
	svcStates   serviceStates
	pendingJobs int
	errMessages []errorMsg
}

func (m model) Init() tea.Cmd {
	return m.initialModel
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case serviceStates:
		m.svcStates = msg

		cmds := make([]tea.Cmd, 0)

		for regName, service := range msg {
			for _, svcName := range maps.Keys(service) {
				regName, svcName := regName, svcName

				cmds = append(cmds, func() tea.Msg {
					return smokeTest(regName, svcName)
				})
			}
		}

		m.pendingJobs = len(cmds)

		return m, tea.Batch(cmds...)
	case error:
		m.errMessages = append(m.errMessages, errorMsg{
			message: msg.Error(),
		})

		return m, tea.Quit
	case serviceStateMsg:
		m.svcStates[msg.region][msg.service] = msg.state
		if msg.err != "" {
			m.errMessages = append(m.errMessages, errorMsg{
				region:  string(msg.region),
				service: string(msg.service),
				message: msg.err,
			})
		}

		m.pendingJobs--

		if m.pendingJobs == 0 {
			return m, tea.Quit
		}

		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	res := ""

	rnames := make([]string, len(m.svcStates))
	for i, r := range maps.Keys(m.svcStates) {
		rnames[i] = string(r)
	}

	slices.Sort(rnames)

	for _, rname := range rnames {
		perState := make(map[string][]serviceName)

		rv := m.svcStates[regionName(rname)]
		for svcName, state := range rv {
			if _, ok := perState[state]; !ok {
				perState[state] = make([]serviceName, 0)
			}

			perState[state] = append(perState[state], svcName)

			slices.Sort(perState[state])
		}

		res += fmt.Sprint(termenv.String(rname + ": ").Foreground(termenv.ANSICyan))

		states := maps.Keys(perState)
		sort.Strings(states)

		for _, state := range states {
			types := perState[state]
			state := termenv.String(state).Foreground(stateColors[state])
			res += fmt.Sprintf(" %v => [ ", state)

			for _, t := range types {
				res += termenv.String(string(t)).Foreground(termenv.ANSIBrightWhite).String() + " "
			}

			res += "] "
		}

		res += "\n"
	}

	for _, msg := range m.errMessages {
		res += termenv.String("ERROR: "+msg.String()).Foreground(termenv.ANSIBrightRed).String() + "\n"
	}

	return res
}

func smokeTestsCommand(cmd *cobra.Command, args []string) {
	mdl := model{svcStates: make(map[regionName]map[serviceName]string)}

	if len(args) == 1 {
		cloud, err := clientconfig.GetCloudFromYAML(&clientconfig.ClientOpts{
			Cloud:    args[0],
			YAMLOpts: config.YAMLOpts{Directory: viper.GetString("os-config-dir")},
		})
		cobra.CheckErr(err)

		mdl.svcStates[regionName(cloud.RegionName)] = nil
	}

	p := tea.NewProgram(mdl)
	if _, err := p.Run(); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Alas, there's been an error: %v", err)

		os.Exit(1)
	}
}

func identityClient(region string) (*gophercloud.ServiceClient, error) {
	opts := &clientconfig.ClientOpts{
		RegionName: region,
		YAMLOpts:   config.YAMLOpts{Directory: viper.GetString("os-config-dir")},
	}

	client, err := clientconfig.NewServiceClient("identity", opts)
	if err != nil {
		return &gophercloud.ServiceClient{}, err
	}

	return client, nil
}

func filterEntries(entries []tokens.CatalogEntry, region regionName) []tokens.CatalogEntry {
	newEntries := []tokens.CatalogEntry{}

	for _, entry := range entries {
		newEndpoints := []tokens.Endpoint{}

		for _, endpoint := range entry.Endpoints {
			if endpoint.Region == string(region) {
				newEndpoints = append(newEndpoints, endpoint)
			}
		}

		entry.Endpoints = newEndpoints

		newEntries = append(newEntries, entry)
	}

	return newEntries
}

func (m model) initialModel() tea.Msg {
	client, err := identityClient("")
	if err != nil {
		return fmt.Errorf("failed to create identity client: %w", err)
	}

	allPages, err := catalog.List(client).AllPages()
	if err != nil {
		return err
	}

	cEntries, err := catalog.ExtractServiceCatalog(allPages)
	if err != nil {
		return err
	}

	newCEntries := cEntries

	if len(m.svcStates) == 1 {
		rname := maps.Keys(m.svcStates)[0]

		newCEntries = filterEntries(cEntries, rname)
	}

	states := make(serviceStates)

	serviceTypesNames := map[serviceType]serviceName{
		"compute":       "nova",
		"identity":      "keystone",
		"image":         "glance",
		"key-manager":   "barbican",
		"load-balancer": "octavia",
		"network":       "neutron",
		"object-store":  "swift",
		"orchestration": "heat",
		"placement":     "placement",
		"volumev3":      "cinder",
	}

	validServiceTypes := maps.Keys(serviceTypesNames)

	for _, entry := range newCEntries {
		stype := serviceType(entry.Type)
		if !slices.Contains(validServiceTypes, stype) {
			continue
		}

		for _, endpoint := range entry.Endpoints {
			rname := regionName(endpoint.Region)

			if _, rok := states[rname]; !rok {
				states[rname] = make(map[serviceName]string)
			}

			states[rname][serviceTypesNames[stype]] = Testing
		}
	}

	return states
}
