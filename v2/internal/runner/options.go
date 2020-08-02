package runner

import (
	"flag"
	"os"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/pkg/requests"
)

// Options contains the configuration options for tuning
// the template requesting process.
type Options struct {
	Debug              bool                   // Debug mode allows debugging request/responses for the engine
	Templates          multiStringFlag        // Signature specifies the template/templates to use
	Target             string                 // Target is a single URL/Domain to scan usng a template
	Targets            string                 // Targets specifies the targets to scan using templates.
	Threads            int                    // Thread controls the number of concurrent requests to make.
	Timeout            int                    // Timeout is the seconds to wait for a response from the server.
	Retries            int                    // Retries is the number of times to retry the request
	Output             string                 // Output is the file to write found subdomains to.
	ProxyURL           string                 // ProxyURL is the URL for the proxy server
	ProxySocksURL      string                 // ProxySocksURL is the URL for the proxy socks server
	Silent             bool                   // Silent suppresses any extra text and only writes found URLs on screen.
	Version            bool                   // Version specifies if we should just show version and exit
	Verbose            bool                   // Verbose flag indicates whether to show verbose output or not
	NoColor            bool                   // No-Color disables the colored output.
	CustomHeaders      requests.CustomHeaders // Custom global headers
	UpdateTemplates    bool                   // UpdateTemplates updates the templates installed at startup
	TemplatesDirectory string                 // TemplatesDirectory is the directory to use for storing templates
	JSON               bool                   // JSON writes json output to files
	JSONRequests       bool                   // write requests/responses for matches in JSON output
	EnableProgressBar  bool                   // Enable progrss bar

	Stdin bool // Stdin specifies whether stdin input was given to the process
}

type multiStringFlag []string

func (m *multiStringFlag) String() string {
	return ""
}

func (m *multiStringFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Target, "target", "", "Target is a single target to scan using template")
	flag.Var(&options.Templates, "t", "Template input file/files to run on host. Can be used multiple times.")
	flag.StringVar(&options.Targets, "l", "", "List of URLs to run templates on")
	flag.StringVar(&options.Output, "o", "", "File to write output to (optional)")
	flag.StringVar(&options.ProxyURL, "proxy-url", "", "URL of the proxy server")
	flag.StringVar(&options.ProxySocksURL, "proxy-socks-url", "", "URL of the proxy socks server")
	flag.BoolVar(&options.Silent, "silent", false, "Show only results in output")
	flag.BoolVar(&options.Version, "version", false, "Show version of nuclei")
	flag.BoolVar(&options.Verbose, "v", false, "Show Verbose output")
	flag.BoolVar(&options.NoColor, "nC", false, "Don't Use colors in output")
	flag.IntVar(&options.Threads, "c", 50, "Number of concurrent requests to make")
	flag.IntVar(&options.Timeout, "timeout", 5, "Time to wait in seconds before timeout")
	flag.IntVar(&options.Retries, "retries", 1, "Number of times to retry a failed request")
	flag.Var(&options.CustomHeaders, "H", "Custom Header.")
	flag.BoolVar(&options.Debug, "debug", false, "Allow debugging of request/responses")
	flag.BoolVar(&options.UpdateTemplates, "update-templates", false, "Update Templates updates the installed templates (optional)")
	flag.StringVar(&options.TemplatesDirectory, "update-directory", "", "Directory to use for storing nuclei-templates")
	flag.BoolVar(&options.JSON, "json", false, "Write json output to files")
	flag.BoolVar(&options.JSONRequests, "json-requests", false, "Write requests/responses for matches in JSON output")
	flag.BoolVar(&options.EnableProgressBar, "pbar", false, "Enable the progress bar")

	flag.Parse()

	// Check if stdin pipe was given
	options.Stdin = hasStdin()

	// Read the inputs and configure the logging
	options.configureOutput()

	// Show the user the banner
	showBanner()

	if options.Version {
		gologger.Infof("Current Version: %s\n", Version)
		os.Exit(0)
	}

	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err := options.validateOptions()
	if err != nil {
		gologger.Fatalf("Program exiting: %s\n", err)
	}
	return options
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}
