package lifecycle

// func (k *Kubernetes) APIVersions(spec types.Spec) []types.Task {
// 	task := plans.StructuredSource{
// 		Producer:  k.producers.APIVersions,
// 		RawPath:   filepath.Join(spec.OutputDir, "api_versions.raw"),
// 		JSONPath:  filepath.Join(spec.OutputDir, "api_versions.json"),
// 		HumanPath: filepath.Join(spec.OutputDir, "api_versions.human"),
// 	}
// 	var err error
// 	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
// 	if err != nil {
// 		task := plans.PreparedError(err, spec)
// 		return []types.Task{task}
// 	}
// 	return []types.Task{&task}
// }
