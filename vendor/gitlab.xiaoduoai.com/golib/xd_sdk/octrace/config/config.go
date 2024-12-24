package config

type Config struct {
	SampleType        int     `json:"sample_type" mapstructure:"sample_type" toml:"sample_type"`
	Fraction          float64 `json:"fraction" mapstructure:"fraction" toml:"fraction"`
	AgentEndPoint     string  `json:"agent_end_point" mapstructure:"agent_end_point" toml:"agent_end_point"`
	CollectorEndPoint string  `json:"collector_end_point" mapstructure:"collector_end_point" toml:"collector_end_point"`
	ServiceName       string  `json:"service_name" mapstructure:"service_name" toml:"service_name"`
	Username          string  `json:"username" mapstructure:"username" toml:"username"`
	Password          string  `json:"password" mapstructure:"password" toml:"password"`
	TraceLogFile      string  `json:"trace_log_file" mapstructure:"trace_log_file" toml:"trace_log_file"`
}

var DefaultConfig = Config{
	SampleType:    1,
	ServiceName:   "unknown",
	AgentEndPoint: "localhost:6831",
}

/* example

if collector end point, will use collector end point, not agent end point.
collector end point like this: "http://localhost:14268/api/traces"
recommend to use local agent



{
	"sample_type"			:1,
	"fraction"				:10.0,
	"agent_end_point"		:"localhost:6831",
	"collector_end_point"	:"",
	"service_name"			:"servicename",
	"username"				:"",
	"password"				:"",
	"trace_log_file"		:"/dev/null"
}

*/
