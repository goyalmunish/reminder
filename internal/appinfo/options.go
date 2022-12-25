package appinfo

type Options struct {
	DataFile string `json:"data_file" yaml:"data_file" mapstructure:"data_file"`
}

func DefaultOptions() *Options {
	dataFilePath := "~/reminder/data.json"
	return &Options{
		DataFile: dataFilePath,
	}
}
