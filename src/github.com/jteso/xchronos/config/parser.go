package config

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	hclObj "github.com/hashicorp/hcl/hcl"
)

func parseVersion(hcltree *hclObj.Object, errors *multierror.Error) string {
	if rawVersion := hcltree.Get("version", false); rawVersion != nil {
		if rawVersion.Len() > 1 {
			errors = multierror.Append(errors, fmt.Errorf("Version was specified more than once"))
		} else {
			if rawVersion.Type != hclObj.ValueTypeString {
				errors = multierror.Append(errors, fmt.Errorf("Version was specified as an invalid type - expected string, found %s", rawVersion.Type))
			} else {
				return rawVersion.Value.(string)
			}
		}
	} else {
		errors = multierror.Append(errors, fmt.Errorf("No version was specified in the configuration"))
	}
	return "*** Error"
}

func parseJobStore(hcltree *hclObj.Object, errors *multierror.Error) []string {
	hosts := make([]string, 0)
	if jobStore := hcltree.Get("job_store", false); jobStore != nil {
		if jobStore.Len() > 1 {
			errors = multierror.Append(errors, fmt.Errorf("job_store was specified more than once"))
		} else {
			if jobStore.Type != hclObj.ValueTypeList {
				errors = multierror.Append(errors, fmt.Errorf("job_store was specified as an invalid type - expected string, found %s", jobStore.Type))
			} else {
				for _, host := range jobStore.Elem(true) {
					hosts = append(hosts, host.Value.(string))
				}
				return hosts
			}
		}
	} else {
		errors = multierror.Append(errors, fmt.Errorf("No job_store was specified in the configuration"))
	}
	return []string{}
}

func parseJobs(hcltree *hclObj.Object, errors *multierror.Error) []JobConfig {
	outputConfig := make([]JobConfig, 0)
	if rawJobs := hcltree.Get("job", false); rawJobs != nil {
		if rawJobs.Type != hclObj.ValueTypeObject {
			errors = multierror.Append(errors, fmt.Errorf("job was specified as an invalid type - expected list, found %s", rawJobs.Type))
		} else {
			arrayJobs := rawJobs.Elem(true)
			for _, job := range arrayJobs {
				outputConfig = append(outputConfig, parseJob(job, errors))
			}
		}

	} else {
		if len(outputConfig) == 0 {
			errors = multierror.Append(errors, fmt.Errorf("No job_store was specified in the configuration"))
		}
	}
	return outputConfig

}

func parseJob(jobNode *hclObj.Object, errors *multierror.Error) JobConfig {
	config := JobConfig{}

	if jobNode.Len() > 1 {
		errors = multierror.Append(errors, fmt.Errorf("job was specified more than once"))
	} else {
		if jobNode.Type != hclObj.ValueTypeObject {
			errors = multierror.Append(errors, fmt.Errorf("job was specified as an invalid type - expected object, found %s", jobNode.Type))
		} else {
			config.Name = jobNode.Key
			config.Trigger = parseJobTrigger(jobNode, errors)
			config.Exec = parseExec(jobNode, errors)
		}
	}
	return config
}

func parseJobTrigger(jobNode *hclObj.Object, errors *multierror.Error) TriggerConfig {
	var triggerConfig TriggerConfig
	if t := jobNode.Get("trigger", false); t != nil {
		err := hcl.DecodeObject(&triggerConfig, t)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return triggerConfig
	// cron := parseCron(jobNode.Get("trigger", false), errors)
	// maxExecs := parseMaxExecutions(jobNode.Get("trigger", false), errors)
	// return TriggerConfig{
	// 	Cron:          cron,
	// 	MaxExecutions: maxExecs,
	// }
}

func parseExec(jobNode *hclObj.Object, errors *multierror.Error) string {
	if rawExec := jobNode.Get("exec", false); rawExec != nil {
		if rawExec.Len() > 1 {
			errors = multierror.Append(errors, fmt.Errorf("exec was specified more than once"))
		} else {
			if rawExec.Type != hclObj.ValueTypeString {
				errors = multierror.Append(errors, fmt.Errorf("exec was specified as an invalid type - expected string, found %s", rawExec.Type))
			} else {
				return rawExec.Value.(string)
			}
		}
	} else {
		errors = multierror.Append(errors, fmt.Errorf("No exec was specified in the configuration"))
	}
	return ""
}

func parseCron(jobNode *hclObj.Object, errors *multierror.Error) string {
	if rawCron := jobNode.Get("cron", false); rawCron != nil {
		if rawCron.Len() > 1 {
			errors = multierror.Append(errors, fmt.Errorf("cron was specified more than once"))
		} else {
			if rawCron.Type != hclObj.ValueTypeString {
				errors = multierror.Append(errors, fmt.Errorf("cron was specified as an invalid type - expected string, found %s", rawCron.Type))
			} else {
				return rawCron.Value.(string)
			}
		}
	} else {
		errors = multierror.Append(errors, fmt.Errorf("No exec was specified in the configuration"))
	}
	return ""
}

func parseMaxExecutions(jobNode *hclObj.Object, errors *multierror.Error) string {
	if rawCron := jobNode.Get("max_executions", false); rawCron != nil {
		if rawCron.Len() > 1 {
			errors = multierror.Append(errors, fmt.Errorf("max_executions was specified more than once"))
		} else {
			if rawCron.Type != hclObj.ValueTypeString {
				errors = multierror.Append(errors, fmt.Errorf("max_executions was specified as an invalid type - expected string, found %s", rawCron.Type))
			} else {
				return rawCron.Value.(string)
			}
		}
	} else {
		errors = multierror.Append(errors, fmt.Errorf("No max_executions was specified in the configuration"))
	}
	return ""
}
