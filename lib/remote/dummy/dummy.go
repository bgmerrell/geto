// Dummy remote for unit testing purposes
package dummy

import (
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/remote"
)

// dummyRemote implements the Remote interface
type dummyRemote struct{}

func New() remote.Remote {
	return new(dummyRemote)
}

func (r dummyRemote) TestConnection(host host.Host) (err error) {
	return nil
}

func (r dummyRemote) Run(host host.Host,
	command string,
	timeout uint32) (stdout string, stderr string, err error) {
	return "test", "", nil
}

func (r dummyRemote) CopyTo(host host.Host,
	recursive bool,
	localPath string,
	remotePath string) (err error) {
	return nil
}

func (r dummyRemote) CopyFrom(host host.Host,
	recursive bool,
	remotePath string,
	localPath string) (err error) {
	return nil
}
