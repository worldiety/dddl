package gen

type Options struct {
	Module              string `json:"module"`
	Dir                 string `json:"dir"`
	DomainPackagePrefix string `json:"domainPackagePrefix"`
}

func DefaultOptions(dir string) Options {
	modname, moddir := ModuleName(dir)
	if moddir == "" {
		moddir = dir
	}

	if modname == "" {
		modname = "example.de/project"
	}

	return Options{
		Module:              modname,
		Dir:                 moddir,
		DomainPackagePrefix: "internal",
	}
}
