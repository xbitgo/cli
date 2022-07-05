package cfg

import (
	jaegerCfg "github.com/uber/jaeger-client-go/config"
)

type Tracing struct {
	ServiceName                string  `json:"service_name" yaml:"service_name"`
	SamplerType                string  `json:"sampler_type" yaml:"sampler_type"`
	SamplerParam               float64 `json:"sampler_param" yaml:"sampler_param"`
	ReporterLocalAgentHostPort string  `json:"reporter_local_agent_host_port" yaml:"reporter_local_agent_host_port"`
	LogSpans                   bool    `json:"log_spans" yaml:"log_spans"`
}

func (c *Tracing) ToJaegerCfg() *jaegerCfg.Configuration {
	samplerType := "const"
	if c.SamplerType != "" {
		samplerType = c.SamplerType
	}
	jc := &jaegerCfg.Configuration{
		ServiceName: c.ServiceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  samplerType,
			Param: c.SamplerParam,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LocalAgentHostPort: c.ReporterLocalAgentHostPort, //
			LogSpans:           c.LogSpans,
		},
	}
	return jc
}
