geto
====

(G)ood (e)nough (t)ask (o)ffloader is a framework for offloading work to hosts
with minimal setup and dependencies.

Basically, geto can be used to offload an arbitrary task to another, target,
host and retrieve results.

You might want to use geto if you have a machine (or more) that needs to
offload work to other machines.  Geto code is only required on machines
from which the work is offloaded; geto is not required (or in any way useful)
on the target host machines.

It is likely that the offloading and result gathering to take on the order of
seconds, so you might not want to use geto if that is a concern.

Here's a trivial example that runs a sleep command on a remote host:

```
package main

import (
    "fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/remote/ssh"
	"github.com/bgmerrell/geto/lib/task"
)

func main() {
	conf, _ := config.ParseConfig("/etc/geto.ini")
	var script task.Script = task.NewScriptWithCommands(
		"sleep", []string{"#!/bin/bash", "sleep 15"}, nil)
	var depFiles []string
	t, _ := task.New(depFiles, script, 0)
	ch := make(chan task.RunOutput)
	go task.RunOnHost(ssh.New(), t, conf.Hosts[0], ch)
	taskOutput := <-ch
	fmt.Printf("stdout: %s\n", taskOutput.Stdout)
	fmt.Printf("stderr: %s\n", taskOutput.Stderr)
	if taskOutput.Err != nil {
		fmt.Printf("err: %s\n", taskOutput.Err)
	}
}
```
## Prerequisites

Any host to which the user wishes to offload must have the following:
* A Unix-like environment (only tested on Linux)
* SSH server with public key authentication with the machine originating the offloading.
* The __timeout__ command in your PATH.  This command is usually installed by default as part of the __coreutils__ package in Linux.

The machine originating the offloading must have the following:
* geto (notice that the target host does not require geto)
* https://github.com/robfig/config
* https://code.google.com/p/go.crypto/
* A Unix-like environment (only tested on Mac OS X)
* SSH client with client key authentication to the target host
* Go (tested on 1.2)

## Terms

* __Host__: Some machine setup with the first set of prerequisites above.
* __Task__: A unit of work to run on a host.  Task IDs are uniquely generated.
* __Script__: A command to run on the host.  The same script may be run by multiple tasks, but a limit can easily be placed on the number of instances of a given script that can be executing concurrently on a given host.

For example, In the above code example, a simple bash script is used to compose a geto script (using the task.NewScriptWithCommands() method).  That geto script is then used to create a new geto task (using the task.New() method).  That task is then executed (using the task.RunOnHost() method) on the first host found in the parsed config file (i.e., conf.Hosts[0]).

## Script details

The Script object consists of a name, commands, and the number of maximum scripts that can run concurrently on a given host.  Otherwise stated:

```
// A script that runs on a target host
type Script struct {
    // Name is the name of a script.  It need not be unique.
	name string
	// The commands that make up a shell-style script.
	// Each index represents a line in the script.
	commands []string
	// The number of scripts of the same name that will run on a target host
	// concurrently.  A nil value means there is no limit.
	maxConcurrent *uint32
}
```

There are multiple ways of creating a script object:

```
func NewScript(name string, maxConcurrent *uint32) Script
```

In the above case the user is responsible for adding the commands to the object.  Alternatively, the commands can be provided when instantiating the script object (which is the strategy used in the first example of this document) like so:

```
func NewScriptWithCommands(name string, commands []string, maxConcurrent *uint32) Script
```

Yet another approach is to provide a path to an existing script file to use to instantiate the geto script:

```
func NewScriptFromPath(name string, path string, maxConcurrent *uint32) (Script, error)
```

Scripts are simply executed on the target host; it is up to the script to indicate how it should be executed (e.g., by using a [shebang interpreter directive](http://en.wikipedia.org/wiki/Shebang_%28Unix%29)).

## Task details

A task object looks like this:
```
// A task that runs on a target host
type Task struct {
    // A unique ID for the task, automatically generated
	Id string
	// A list of files and/or directories that the task requires
	DepFiles []string
	// A script for the task to run
	Script Script
	// The number of seconds before giving up on a task after it has been
	// started
	Timeout uint32
}
```

Any file dependencies (specified by DepFiles) are copied to the target host and placed in a special "DEPS" directory.  The script is also copied to the target host and placed in the same parent directory as the "DEPS" directory.  This means that file dependencies can be relatively referenced from the script.  For example, a foo.bin file dependency could be referenced in the script by "DEPS/foo.bin".  (NOTE: This may or may not be tested at this point).

There is currently one way to instantiate a task object:

```
func New(depFiles []string, script Script, timeout uint32) (Task, error)
```

Once a task has been created, however, there are several fun ways to run it.  The user can provide exactly which host on which the task should be run, like this:

```
func RunOnHost(conn remote.Remote, task Task, host host.Host, resultChan chan<- RunOutput)
```

Or, the user might wish to just have a random host picked, like this:

```
func RunOnRandomHost(conn remote.Remote, task Task, ch chan<- RunOutput)
```

The user can also perform basic load balancing by having geto choose the host that is running the fewest instances of a task's script, like this:

```
func RunOnHostBalancedByScriptName(conn remote.Remote, task Task, ch chan<- RunOutput)
```

## TODO

* Allow the remote copy operations to be done using password authentication (see [issue #1](https://github.com/bgmerrell/geto/issues/1))
* Implement Python bridge allowing geto to be wielded from Python.  There is already a proof-of-concept code checked into the geto repo.  The code consists of a Go JSON rpc server and Python RPC client that calls it.
* Various TODO-marked code.
