package govalid

// IntPort returns an integer validation rule for port.
func IntPort() Rule { return IntBetween(1, 65535) }

// IntPercentage returns an integer validation rule for percentage.
func IntPercentage() Rule { return IntBetween(0, 100) }

// IntHTTPStatus returns an integer validation rule for http status.
func IntHTTPStatus() Rule { return IntBetween(100, 599) }
