package ssh

import (
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/remote"
	"github.com/bgmerrell/geto/lib/ssh"
)

// remote implements the Remote interface
type sshRemote struct{}

func New() remote.Remote {
	return new(sshRemote)
}

func (r sshRemote) TestConnection(host host.Host) (err error) {
	conf := config.GetParsedConfig()
	return ssh.TestConnection(
		host.Addr,
		host.Username,
		host.Password,
		conf.PrivKeyPath,
		host.PortNum)
}

func (r sshRemote) Run(host host.Host,
	command string,
	timeout uint32) (stdout string, stderr string, err error) {
	conf := config.GetParsedConfig()
	return ssh.Run(
		host.Addr,
		host.Username,
		host.Password,
		conf.PrivKeyPath,
		host.PortNum,
		command,
		timeout)
}

func (r sshRemote) CopyTo(host host.Host,
	recursive bool,
	localPath string,
	remotePath string) (err error) {
	return ssh.ScpTo(
		host.Addr,
		host.Username,
		host.PortNum,
		recursive,
		localPath,
		remotePath)
}

func (r sshRemote) CopyFrom(host host.Host,
	recursive bool,
	remotePath string,
	localPath string) (err error) {
	return ssh.ScpFrom(
		host.Addr,
		host.Username,
		host.PortNum,
		recursive,
		remotePath,
		localPath)
}
