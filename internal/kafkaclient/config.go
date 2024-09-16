package kafkaclient

type Config struct {
	AutoCommit         bool     `json:"auto_commit"           koanf:"auto_commit"`
	AutoOffsetReset    string   `json:"auto_offset_reset"     koanf:"auto_offset_reset"`
	BootstrapServers   []string `json:"bootstrap_servers"     koanf:"bootstrap_servers"`
	Debug              bool     `json:"debug"                 koanf:"debug"`
	GroupId            string   `json:"group_id"              koanf:"group_id"`
	HealthcheckTimeout int32    `json:"healthcheck_timeout"   koanf:"healthcheck_timeout"`
	SaslMechanisms     string   `json:"sasl_mechanisms"       koanf:"sasl_mechanisms"`
	SaslPassword       string   `json:"sasl_password"         koanf:"sasl_password"`
	SaslUsername       string   `json:"sasl_username"         koanf:"sasl_username"`
	SecurityProtocol   string   `json:"security_protocol"     koanf:"security_protocol"`
	TimestampMode      string   `json:"timestamp_mode"        koanf:"timestamp_mode"`
	Topics             []string `json:"topic"                 koanf:"topic"`
}
