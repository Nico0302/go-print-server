package printer

type JobAttributes struct {
	Media                string   `mapstructure:"media,omitempty"`
	MediaCol             MediaCol `mapstructure:"media-col,omitempty"`
	OrientationRequested int      `mapstructure:"orientation-requested,omitempty"`
	PrintQuality         int      `mapstructure:"print-quality,omitempty"`
	PrintColorMode       string   `mapstructure:"print-color-mode,omitempty"`
	Sides                string   `mapstructure:"sides,omitempty"`
}

type MediaCol struct {
	MediaSource   string `mapstructure:"media-source,omitempty"`
	MediaType     string `mapstructure:"media-type,omitempty"`
	MediaSizeName string `mapstructure:"media-size-name,omitempty"`
}
