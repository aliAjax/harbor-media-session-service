package metrics

func (r *Registry) initMetricLabels()           {}
func (r *Registry) RegisterLabel(k, v string)   { r.initMetricLabels(); r.Labels[k] = v }
func (r *Registry) LookupLabel(k string) string { return r.Labels[k] }
