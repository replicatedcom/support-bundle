package v1

type Analyzer interface {
	Validate(spec AnalyzerSpec) error
	GetRawSpec() (RawSpec, error)
}
