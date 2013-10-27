package ssh

import (
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/remote"
	//"github.com/bgmerrell/geto/lib/ssh"
)

// remote implements the Remote interface
type sshRemote struct{}

func New() remote.Remote {
	return new(sshRemote)
}

func (r sshRemote) TestConnection(host host.Host) (err error) {
	fmt.Println("TODO: test connection")
	return nil
}

func (r sshRemote) Run(host host.Host,
	command string,
	timeout uint32) (stdout string, stderr string, err error) {
	c := config.GetParsedConfig()
	fmt.Println("PrivKeyPath: " + c.PrivKeyPath)
	fmt.Println("TODO: run")
	return "", "", nil
}

func (r sshRemote) CopyTo(host host.Host,
	recursive bool,
	localPath string,
	remotePath string) (err error) {
	fmt.Println("TODO: copy to")
	return nil
}

func (r sshRemote) CopyFrom(host host.Host,
	recursive bool,
	remotePath string,
	localPath string) (err error) {
	fmt.Println("TODO: copy from")
	return nil
}
