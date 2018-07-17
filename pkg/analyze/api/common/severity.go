package common

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
	SeverityInfo  Severity = "info"
)

type Severity string

func SeverityOrDefault(severity Severity) Severity {
	if severity != "" {
		return severity
	}
	return SeverityError
}

func SeverityCompare(severity, threshold Severity) int {
	switch threshold {
	case SeverityError:
		switch severity {
		case SeverityError:
			return 0
		default:
			return -1
		}
	case SeverityWarn:
		switch severity {
		case SeverityError:
			return 1
		case SeverityWarn:
			return 0
		default:
			return -1
		}
	case SeverityInfo:
		switch severity {
		case SeverityError, SeverityWarn:
			return 1
		case SeverityInfo:
			return 0
		default:
			return -1
		}
	default:
		switch severity {
		case SeverityError, SeverityWarn, SeverityInfo:
			return 1
		default:
			return 0
		}
	}
}
