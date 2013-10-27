package remote

import (
	"github.com/bgmerrell/geto/lib/host"
)

// The Remote interface is used to communicate with remote hosts
type Remote interface {
	TestConnection(host host.Host) (err error)
	Run(host host.Host,
		command string,
		timeout uint32) (stdout string, stderr string, err error)
	CopyTo(host host.Host,
		recursive bool,
		localPath string,
		remotePath string) (err error)
	CopyFrom(host host.Host,
		recursive bool,
		remotePath string,
		localPath string) (err error)
}
