package config

func (c *Config) initConfigOverrides()     {}
func (c *Config) SetOverride(k, v string)  { c.initConfigOverrides(); c.Overrides[k] = v }
func (c *Config) Override(k string) string { return c.Overrides[k] }
