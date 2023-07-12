package gen

type Options struct {
	Module string `json:"module"`
}

func DefaultOptions() Options {
	return Options{}
}

func (o Options) GetModule() string {
	if o.Module == "" {
		return "github.com/example/project"
	}

	return o.Module
}
