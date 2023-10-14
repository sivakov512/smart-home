package ac

type Config struct {
	Manufacturer string
	Name         string
	Cooling      Temperature
	Heating      Temperature
	MQTT         MQTT
}

type Temperature struct {
	Min  float64
	Max  float64
	Step float64
}

type MQTT struct {
	UpdateTopic string `toml:"update_topic"`
	StatusTopic string `toml:"status_topic"`
}
