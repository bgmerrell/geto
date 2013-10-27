package remote

import (
	"github.com/bgmerrell/geto/lib/host"
)

// The Remote interface is used to communicate with remote hosts
type Remote interface {
	// TestConnection returns an error if the host can't be communicated
	// with.
	TestConnection(host host.Host) (err error)

	// Run runs command on the host with the specified timeout (in seconds.
	//
	// stdout and stderr strings from running the remote command are
	// returned.  
	//
	// The returned error is nil if the command runs, has no problems
	// copying stdin, stdout, and stderr, and exits with a zero exit
	// status.
	Run(host host.Host,
		command string,
		timeout uint32) (stdout string, stderr string, err error)

	// CopyTo copies the local localPath to remotePath on host.
	//
	// The copy is performed recursively if recursive is true.
	//
	// An error is returned if the copy fails.
	CopyTo(host host.Host,
		recursive bool,
		localPath string,
		remotePath string) (err error)

	// CopyFrom copies from remotePath on host to localPath locally.
	//
	// The copy is performed recursively if recursive is true.
	//
	// An error is returned if the copy fails.
	CopyFrom(host host.Host,
		recursive bool,
		remotePath string,
		localPath string) (err error)
}
