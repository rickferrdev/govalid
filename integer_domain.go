package govalid

func IntPort() Rule       { return IntBetween(1, 65535) }
func IntPercentage() Rule { return IntBetween(0, 100) }
func IntHTTPStatus() Rule { return IntBetween(100, 599) }
