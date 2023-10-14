package common

type TemperatureConfig struct {
	Min  float64
	Max  float64
	Step float64
}

type MQTTConfig struct {
	UpdateTopic string `toml:"update_topic"`
	StatusTopic string `toml:"status_topic"`
}
