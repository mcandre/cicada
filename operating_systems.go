package cicada

// OperatingSystems enumerates OS endoflife.date products.
var OperatingSystems = []string{
	"almalinux",
	"alpine",
	"amazon-linux",
	"android",
	"centos",
	"debian",
	"fedora",
	"freebsd",
	"iphone",
	"kindle",
	"linux",
	"linuxmint",
	"macos",
	"nixos",
	"openbsd",
	"opensuse",
	"rhel",
	"rocky-linux",
	"ros",
	"ubuntu",
	"windows",
	"windowsembedded",
	"windowsserver",
	"yocto",
}

// IsOperatingSystem returns true for known OS endoflife.date products.
//
// Otherwise, returns false.
func IsOperatingSystem(product string) bool {
	for _, operatingSystem := range OperatingSystems {
		if product == operatingSystem {
			return true
		}
	}

	return false
}
