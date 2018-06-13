package graphkit

// ValidationConfig is used to control which errors are ignored / raised / corrected by the validator.
type ValidationConfig struct {
	IgnoreNonExistentTypes   bool
	IgnoreExponentialQueries bool
}
