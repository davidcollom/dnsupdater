package api

type ProviderConfig struct {
	Type string                 `mapstructure:"type"`
	Opts map[string]interface{} `mapstructure:"opts"`
}
