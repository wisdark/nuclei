package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/internal/runner"
	"github.com/projectdiscovery/nuclei/v2/pkg/model/types/severity"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

var (
	cfgFile string
	options = &types.Options{}
)

func main() {
	readConfig()

	runner.ParseOptions(options)

	nucleiRunner, err := runner.New(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}
	if nucleiRunner == nil {
		return
	}
	if err := nucleiRunner.RunEnumeration(); err != nil {
		gologger.Fatal().Msgf("Could not run nuclei: %s\n", err)
	}
	nucleiRunner.Close()
}

func readConfig() {
	home, _ := os.UserHomeDir()
	templatesDirectory := filepath.Join(home, "nuclei-templates")

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Nuclei is a fast, template based vulnerability scanner focusing
on extensive configurability, massive extensibility and ease of use.`)

	/* TODO Important: The defined default values, especially for slice/array types are NOT DEFAULT VALUES, but rather implicit values to which the user input is appended.
	This can be very confusing and should be addressed
	*/

	createGroup(flagSet, "input", "Target",
		flagSet.StringSliceVarP(&options.Targets, "target", "u", []string{}, "target URLs/hosts to scan"),
		flagSet.StringVarP(&options.TargetsFilePath, "list", "l", "", "path to file containing a list of target URLs/hosts to scan (one per line)"),
	)

	createGroup(flagSet, "templates", "Templates",
		flagSet.StringSliceVarP(&options.Templates, "templates", "t", []string{}, "template or template directory paths to include in the scan"),
		flagSet.BoolVarP(&options.NewTemplates, "new-templates", "nt", false, "run only new templates added in latest nuclei-templates release"),
		flagSet.StringSliceVarP(&options.Workflows, "workflows", "w", []string{}, "workflow or workflow directory paths to include in the scan"),
		flagSet.BoolVar(&options.Validate, "validate", false, "validate the passed templates to nuclei"),
		flagSet.BoolVar(&options.TemplateList, "tl", false, "list all available templates"),
	)

	createGroup(flagSet, "filters", "Filtering",
		flagSet.NormalizedStringSliceVar(&options.Tags, "tags", []string{}, "execute a subset of templates that contain the provided tags"),
		flagSet.NormalizedStringSliceVarP(&options.IncludeTags, "include-tags", "itags", []string{}, "tags from the default deny list that permit executing more intrusive templates"), // TODO show default deny list
		flagSet.NormalizedStringSliceVarP(&options.ExcludeTags, "exclude-tags", "etags", []string{}, "exclude templates with the provided tags"),
		flagSet.StringSliceVarP(&options.IncludeTemplates, "include-templates", "it", []string{}, "templates to be executed even if they are excluded either by default or configuration"),
		flagSet.StringSliceVarP(&options.ExcludedTemplates, "exclude-templates", "et", []string{}, "template or template directory paths to exclude"),
		flagSet.VarP(&options.Severities, "severity", "s", fmt.Sprintf("Templates to run based on severity. Possible values: %s", severity.GetSupportedSeverities().String())),
		flagSet.VarP(&options.ExcludeSeverities, "exclude-severity", "es", fmt.Sprintf("Templates to exclude based on severity. Possible values: %s", severity.GetSupportedSeverities().String())),
		flagSet.NormalizedStringSliceVarP(&options.Author, "author", "a", []string{}, "execute templates that are (co-)created by the specified authors"),
	)

	createGroup(flagSet, "output", "Output",
		flagSet.StringVarP(&options.Output, "output", "o", "", "output file to write found issues/vulnerabilities"),
		flagSet.BoolVar(&options.Silent, "silent", false, "display findings only"),
		flagSet.BoolVarP(&options.NoColor, "no-color", "nc", false, "disable output content coloring (ANSI escape codes)"),
		flagSet.BoolVar(&options.JSON, "json", false, "write output in JSONL(ines) format"),
		flagSet.BoolVarP(&options.JSONRequests, "include-rr", "irr", false, "include request/response pairs in the JSONL output (for findings only)"),
		flagSet.BoolVarP(&options.NoMeta, "no-meta", "nm", false, "don't display match metadata"),
		flagSet.BoolVarP(&options.NoTimestamp, "no-timestamp", "nts", false, "don't display timestamp metadata in CLI output"),
		flagSet.StringVarP(&options.ReportingDB, "report-db", "rdb", "", "local nuclei reporting database (always use this to persist report data)"),
		flagSet.StringVarP(&options.MarkdownExportDirectory, "markdown-export", "me", "", "directory to export results in markdown format"),
		flagSet.StringVarP(&options.SarifExport, "sarif-export", "se", "", "file to export results in SARIF format"),
	)

	createGroup(flagSet, "configs", "Configurations",
		flagSet.StringVar(&cfgFile, "config", "", "path to the nuclei configuration file"),
		flagSet.StringVarP(&options.ReportingConfig, "report-config", "rc", "", "nuclei reporting module configuration file"), // TODO merge into the config file or rename to issue-tracking
		flagSet.StringSliceVarP(&options.CustomHeaders, "header", "H", []string{}, "custom headers in header:value format"),
		flagSet.RuntimeMapVarP(&options.Vars, "var", "V", []string{}, "custom vars in var=value format"),
		flagSet.StringVarP(&options.ResolversFile, "resolvers", "r", "", "file containing resolver list for nuclei"),
		flagSet.BoolVarP(&options.SystemResolvers, "system-resolvers", "sr", false, "use system DNS resolving as error fallback"),
		flagSet.BoolVar(&options.OfflineHTTP, "passive", false, "enable passive HTTP response processing mode"),
		flagSet.BoolVarP(&options.EnvironmentVariables, "env-vars", "ev", false, "enable environment variables to be used in template"),
	)

	createGroup(flagSet, "interactsh", "interactsh",
		flagSet.StringVarP(&options.InteractshURL, "interactsh-server", "iserver", "https://interactsh.com", "interactsh server url for self-hosted instance"),
		flagSet.StringVarP(&options.InteractshToken, "interactsh-token", "itoken", "", "authentication token for self-hosted interactsh server"),
		flagSet.IntVar(&options.InteractionsCacheSize, "interactions-cache-size", 5000, "number of requests to keep in the interactions cache"),
		flagSet.IntVar(&options.InteractionsEviction, "interactions-eviction", 60, "number of seconds to wait before evicting requests from cache"),
		flagSet.IntVar(&options.InteractionsPollDuration, "interactions-poll-duration", 5, "number of seconds to wait before each interaction poll request"),
		flagSet.IntVar(&options.InteractionsColldownPeriod, "interactions-cooldown-period", 5, "extra time for interaction polling before exiting"),
		flagSet.BoolVarP(&options.NoInteractsh, "no-interactsh", "ni", false, "disable interactsh server for OAST testing, exclude OAST based templates"),
	)

	createGroup(flagSet, "rate-limit", "Rate-Limit",
		flagSet.IntVarP(&options.RateLimit, "rate-limit", "rl", 150, "maximum number of requests to send per second"),
		flagSet.IntVarP(&options.RateLimitMinute, "rate-limit-minute", "rlm", 0, "maximum number of requests to send per minute"),
		flagSet.IntVarP(&options.BulkSize, "bulk-size", "bs", 25, "maximum number of hosts to be analyzed in parallel per template"),
		flagSet.IntVarP(&options.TemplateThreads, "concurrency", "c", 25, "maximum number of templates to be executed in parallel"),
	)

	createGroup(flagSet, "optimization", "Optimizations",
		flagSet.IntVar(&options.Timeout, "timeout", 5, "time to wait in seconds before timeout"),
		flagSet.IntVar(&options.Retries, "retries", 1, "number of times to retry a failed request"),
		flagSet.IntVarP(&options.MaxHostError, "max-host-error", "mhe", 30, "max errors for a host before skipping from scan"),
		flagSet.BoolVar(&options.Project, "project", false, "use a project folder to avoid sending same request multiple times"),
		flagSet.StringVar(&options.ProjectPath, "project-path", os.TempDir(), "set a specific project path"),
		flagSet.BoolVarP(&options.StopAtFirstMatch, "stop-at-first-path", "spm", false, "stop processing HTTP requests after the first match (may break template/workflow logic)"),
		flagSet.BoolVar(&options.Stream, "stream", false, "Stream mode - start elaborating without sorting the input"),
	)

	createGroup(flagSet, "headless", "Headless",
		flagSet.BoolVar(&options.Headless, "headless", false, "enable templates that require headless browser support"),
		flagSet.IntVar(&options.PageTimeout, "page-timeout", 20, "seconds to wait for each page in headless mode"),
		flagSet.BoolVarP(&options.ShowBrowser, "show-browser", "sb", false, "show the browser on the screen when running templates with headless mode"),
		flagSet.BoolVarP(&options.UseInstalledChrome, "system-chrome", "sc", false, "Use local installed chrome browser instead of nuclei installed"),
	)

	createGroup(flagSet, "debug", "Debug",
		flagSet.BoolVar(&options.Debug, "debug", false, "show all requests and responses"),
		flagSet.BoolVar(&options.DebugRequests, "debug-req", false, "show all sent requests"),
		flagSet.BoolVar(&options.DebugResponse, "debug-resp", false, "show all received responses"),

		/* TODO why the separation? http://proxy:port vs socks5://proxy:port etc
		   TODO should auto-set the HTTP_PROXY variable for the process? */
		flagSet.StringVarP(&options.ProxyURL, "proxy-url", "proxy", "", "URL of the HTTP proxy server"),
		flagSet.StringVar(&options.ProxySocksURL, "proxy-socks-url", "", "URL of the SOCKS proxy server"),
		flagSet.StringVarP(&options.TraceLogFile, "trace-log", "tlog", "", "file to write sent requests trace log"),
		flagSet.BoolVar(&options.Version, "version", false, "show nuclei version"),
		flagSet.BoolVarP(&options.Verbose, "verbose", "v", false, "show verbose output"),
		flagSet.BoolVar(&options.VerboseVerbose, "vv", false, "display templates loaded for scan"),
		flagSet.BoolVarP(&options.TemplatesVersion, "templates-version", "tv", false, "shows the version of the installed nuclei-templates"),
	)

	createGroup(flagSet, "update", "Update",
		flagSet.BoolVar(&options.UpdateNuclei, "update", false, "update nuclei engine to the latest released version"),
		flagSet.BoolVarP(&options.UpdateTemplates, "update-templates", "ut", false, "update nuclei-templates to latest released version"),
		flagSet.StringVarP(&options.TemplatesDirectory, "update-directory", "ud", templatesDirectory, "overwrite the default directory to install nuclei-templates"),
		flagSet.BoolVarP(&options.NoUpdateTemplates, "disable-update-check", "duc", false, "disable automatic nuclei/templates update check"),
	)

	createGroup(flagSet, "stats", "Statistics",
		flagSet.BoolVar(&options.EnableProgressBar, "stats", false, "display statistics about the running scan"),
		flagSet.BoolVarP(&options.StatsJSON, "stats-json", "sj", false, "write statistics data to an output file in JSONL(ines) format"),
		flagSet.IntVarP(&options.StatsInterval, "stats-interval", "si", 5, "number of seconds to wait between showing a statistics update"),
		flagSet.BoolVarP(&options.Metrics, "metrics", "m", false, "expose nuclei metrics on a port"),
		flagSet.IntVarP(&options.MetricsPort, "metrics-port", "mp", 9092, "port to expose nuclei metrics on"),
	)

	_ = flagSet.Parse()

	if cfgFile != "" {
		if err := flagSet.MergeConfigFile(cfgFile); err != nil {
			gologger.Fatal().Msgf("Could not read config: %s\n", err)
		}
	}
}

func createGroup(flagSet *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

/*
HacktoberFest update: Below, you can find our ticket recommendations. Tasks with the "good first issue" label are suitable for first time contributors. If you have other ideas, or need help with getting started, join our Discord channel or reach out to @forgedhallpass.

https://github.com/issues?q=is%3Aopen+is%3Aissue+user%3Aprojectdiscovery+label%3AHacktoberfest

*/
