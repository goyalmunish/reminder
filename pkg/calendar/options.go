package calendar

type Options struct {
	CredentialFile string `json:"credential_file" yaml:"credential_file" mapstructure:"credential_file"`
}

func DefaultOptions() *Options {
	return &Options{
		CredentialFile: "~/calendar_credentials.json",
	}
}
