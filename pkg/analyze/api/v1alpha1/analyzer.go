package v1alpha1

type Analyzer interface {
	Validate(spec AnalyzerSpec) error
	GetRawSpec() (RawSpec, error)
}
