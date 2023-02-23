package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/projectdiscovery/nuclei/v2/pkg/testutils"
)

var (
	debug        = os.Getenv("DEBUG") == "true"
	githubAction = os.Getenv("GH_ACTION") == "true"
	customTests  = os.Getenv("TESTS")

	success = aurora.Green("[✓]").String()
	failed  = aurora.Red("[✘]").String()

	protocolTests = map[string]map[string]testutils.TestCase{
		"http":            httpTestcases,
		"network":         networkTestcases,
		"dns":             dnsTestCases,
		"workflow":        workflowTestcases,
		"loader":          loaderTestcases,
		"websocket":       websocketTestCases,
		"headless":        headlessTestcases,
		"whois":           whoisTestCases,
		"ssl":             sslTestcases,
		"code":            codeTestcases,
		"templatesPath":   templatesPathTestCases,
		"templatesDir":    templatesDirTestCases,
		"file":            fileTestcases,
		"offlineHttp":     offlineHttpTestcases,
		"customConfigDir": customConfigDirTestCases,
	}

	// For debug purposes
	runProtocol = ""
	runTemplate = ""
	extraArgs   = []string{}
)

func main() {
	flag.StringVar(&runProtocol, "protocol", "", "run integration tests of given protocol")
	flag.StringVar(&runTemplate, "template", "", "run integration test of given template")
	flag.Parse()

	// allows passing extra args to nuclei
	eargs := os.Getenv("DebugExtraArgs")
	if eargs != "" {
		extraArgs = strings.Split(eargs, " ")
		testutils.ExtraDebugArgs = extraArgs
	}

	if runProtocol != "" {
		debug = true
		debugTests()
		os.Exit(1)
	}

	failedTestTemplatePaths := runTests(toMap(toSlice(customTests)))

	if len(failedTestTemplatePaths) > 0 {
		if githubAction {
			debug = true
			fmt.Println("::group::Failed integration tests in debug mode")
			_ = runTests(failedTestTemplatePaths)
			fmt.Println("::endgroup::")
		}

		os.Exit(1)
	}
}

func debugTests() {
	for tpath, testcase := range protocolTests[runProtocol] {
		if runTemplate != "" && !strings.Contains(tpath, runTemplate) {
			continue
		}
		if err := testcase.Execute(tpath); err != nil {
			fmt.Printf("\n%v", err.Error())
		}
	}
}

func runTests(customTemplatePaths map[string]struct{}) map[string]struct{} {
	failedTestTemplatePaths := map[string]struct{}{}

	for proto, testCases := range protocolTests {
		if len(customTemplatePaths) == 0 {
			fmt.Printf("Running test cases for %q protocol\n", aurora.Blue(proto))
		}

		for templatePath, testCase := range testCases {
			if len(customTemplatePaths) == 0 || contains(customTemplatePaths, templatePath) {
				if failedTemplatePath, err := execute(testCase, templatePath); err != nil {
					failedTestTemplatePaths[failedTemplatePath] = struct{}{}
				}
			}
		}
	}

	return failedTestTemplatePaths
}

func execute(testCase testutils.TestCase, templatePath string) (string, error) {
	if err := testCase.Execute(templatePath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s Test \"%s\" failed: %s\n", failed, templatePath, err)
		return templatePath, err
	}

	fmt.Printf("%s Test \"%s\" passed!\n", success, templatePath)
	return "", nil
}

func expectResultsCount(results []string, expectedNumber int) error {
	if len(results) != expectedNumber {
		return fmt.Errorf("incorrect number of results: %d (actual) vs %d (expected) \nResults:\n\t%s\n", len(results), expectedNumber, strings.Join(results, "\n\t"))
	}
	return nil
}

func toSlice(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}

	return strings.Split(value, ",")
}

func toMap(slice []string) map[string]struct{} {
	result := make(map[string]struct{}, len(slice))
	for _, value := range slice {
		if _, ok := result[value]; !ok {
			result[value] = struct{}{}
		}
	}
	return result
}

func contains(input map[string]struct{}, value string) bool {
	_, ok := input[value]
	return ok
}
