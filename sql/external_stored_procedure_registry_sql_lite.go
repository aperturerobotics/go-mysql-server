//go:build sql_lite

package sql

// ExternalStoredProcedureRegistry omits reflected procedure dispatch in lite builds.
type ExternalStoredProcedureRegistry struct{}

// NewExternalStoredProcedureRegistry returns the lite registry.
func NewExternalStoredProcedureRegistry() ExternalStoredProcedureRegistry {
	return ExternalStoredProcedureRegistry{}
}

// Register omits external procedures in lite builds.
func (*ExternalStoredProcedureRegistry) Register(ExternalStoredProcedureDetails) {}

// LookupByName returns no external procedures in lite builds.
func (*ExternalStoredProcedureRegistry) LookupByName(string) ([]ExternalStoredProcedureDetails, error) {
	return nil, nil
}

// LookupByNameAndParamCount returns no external procedures in lite builds.
func (*ExternalStoredProcedureRegistry) LookupByNameAndParamCount(string, int) (*ExternalStoredProcedureDetails, error) {
	return nil, nil
}
