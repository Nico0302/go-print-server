package printer

type JobAttributes struct {
	MediaCol MediaCol `mapstructure:"media-col"`
}

type MediaCol struct {
	MediaSource string `mapstructure:"media-source"`
	MediaType   string `mapstructure:"media-type"`
	Media       string `mapstructure:"media"`
}
