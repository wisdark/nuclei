package runner

import (
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/pkg/templates"
	"github.com/remeh/sizedwaitgroup"
	"go.uber.org/atomic"
)

// processSelfContainedTemplates execute a self-contained template.
func (r *Runner) processSelfContainedTemplates(template *templates.Template) bool {
	match, err := template.Executer.Execute("")
	if err != nil {
		gologger.Warning().Msgf("[%s] Could not execute step: %s\n", r.colorizer.BrightBlue(template.ID), err)
	}
	return match
}

// processTemplateWithList execute a template against the list of user provided targets
func (r *Runner) processTemplateWithList(template *templates.Template) bool {
	results := &atomic.Bool{}
	wg := sizedwaitgroup.New(r.options.BulkSize)
	processItem := func(k, _ []byte) error {
		URL := string(k)

		// Skip if the host has had errors
		if r.hostErrors != nil && r.hostErrors.Check(URL) {
			return nil
		}
		wg.Add()
		go func(URL string) {
			defer wg.Done()

			match, err := template.Executer.Execute(URL)
			if err != nil {
				gologger.Warning().Msgf("[%s] Could not execute step: %s\n", r.colorizer.BrightBlue(template.ID), err)
			}
			results.CAS(false, match)
		}(URL)
		return nil
	}
	if r.options.Stream {
		_ = r.hostMapStream.Scan(processItem)
	} else {
		r.hostMap.Scan(processItem)
	}

	wg.Wait()
	return results.Load()
}

// processTemplateWithList process a template on the URL list
func (r *Runner) processWorkflowWithList(template *templates.Template) bool {
	results := &atomic.Bool{}
	wg := sizedwaitgroup.New(r.options.BulkSize)

	processItem := func(k, _ []byte) error {
		URL := string(k)

		// Skip if the host has had errors
		if r.hostErrors != nil && r.hostErrors.Check(URL) {
			return nil
		}
		wg.Add()
		go func(URL string) {
			defer wg.Done()
			match := template.CompiledWorkflow.RunWorkflow(URL)
			results.CAS(false, match)
		}(URL)
		return nil
	}

	if r.options.Stream {
		_ = r.hostMapStream.Scan(processItem)
	} else {
		r.hostMap.Scan(processItem)
	}

	wg.Wait()
	return results.Load()
}
