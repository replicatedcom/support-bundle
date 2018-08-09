package v1

type Requirement interface {
	Validate(spec AnalyzerSpec) error
	GetRawSpec() (RawSpec, error)
}
