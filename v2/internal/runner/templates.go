package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/karrick/godirwalk"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/pkg/model/types/severity"
	"github.com/projectdiscovery/nuclei/v2/pkg/parsers"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

func (r *Runner) templateLogMsg(id, name, author string, templateSeverity severity.Severity) string {
	// Display the message for the template
	return fmt.Sprintf("[%s] %s (%s) [%s]",
		r.colorizer.BrightBlue(id).String(),
		r.colorizer.Bold(name).String(),
		r.colorizer.BrightYellow(appendAtSignToAuthors(author)).String(),
		r.addColor(templateSeverity))
}

// appendAtSignToAuthors appends @ before each author and returns final string
func appendAtSignToAuthors(author string) string {
	authors := strings.Split(author, ",")
	if len(authors) == 0 {
		return "@none"
	}
	if len(authors) == 1 {
		if !strings.HasPrefix(authors[0], "@") {
			return fmt.Sprintf("@%s", authors[0])
		}
		return authors[0]
	}
	values := make([]string, 0, len(authors))
	for _, k := range authors {
		if !strings.HasPrefix(authors[0], "@") {
			values = append(values, fmt.Sprintf("@%s", k))
		} else {
			values = append(values, k)
		}
	}
	return strings.Join(values, ",")
}

func (r *Runner) logAvailableTemplate(tplPath string) {
	t, err := parsers.ParseTemplate(tplPath)
	if err != nil {
		gologger.Error().Msgf("Could not parse file '%s': %s\n", tplPath, err)
	} else {
		gologger.Print().Msgf("%s\n", r.templateLogMsg(t.ID,
			types.ToString(t.Info.Name),
			types.ToString(t.Info.Authors),
			t.Info.SeverityHolder.Severity))
	}
}

// ListAvailableTemplates prints available templates to stdout
func (r *Runner) listAvailableTemplates() {
	if r.templatesConfig == nil {
		return
	}

	if _, err := os.Stat(r.templatesConfig.TemplatesDirectory); os.IsNotExist(err) {
		gologger.Error().Msgf("%s does not exists", r.templatesConfig.TemplatesDirectory)
		return
	}

	gologger.Print().Msgf(
		"\nListing available v.%s nuclei templates for %s",
		r.templatesConfig.CurrentVersion,
		r.templatesConfig.TemplatesDirectory,
	)
	err := directoryWalker(
		r.templatesConfig.TemplatesDirectory,
		func(path string, d *godirwalk.Dirent) error {
			if d.IsDir() && path != r.templatesConfig.TemplatesDirectory {
				gologger.Print().Msgf("\n%s:\n\n", r.colorizer.Bold(r.colorizer.BgBrightBlue(d.Name())).String())
			} else if strings.HasSuffix(path, ".yaml") {
				r.logAvailableTemplate(path)
			}
			return nil
		},
	)
	// directory couldn't be walked
	if err != nil {
		gologger.Error().Msgf("Could not find templates in directory '%s': %s\n", r.templatesConfig.TemplatesDirectory, err)
	}
}

func directoryWalker(fsPath string, callback func(fsPath string, d *godirwalk.Dirent) error) error {
	return godirwalk.Walk(fsPath, &godirwalk.Options{
		Callback: callback,
		ErrorCallback: func(fsPath string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
}
