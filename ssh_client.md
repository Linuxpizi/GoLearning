# SSH Client connection in Golang

SSH is a network protocol for establishing a secure shell session on distant servers. In Golang the package godoc.org/golang.org/x/crypto/ssh implements SSH client and SSH server.

In this article, we are using SSH client to run a shell command on a remote machine. Every SSH connection requires an ssh.CleintConfig object that defines configuration options such as authentication.

Authentication Options
Depending on how the remote server is configure, there are two ways to authenticate:

using a username and SSH certificate
using a username and password credentials
If you want to authenticate with username and password you should create ssh.ClientConfig in the following way:

```golang
sshConfig := &ssh.ClientConfig{
	User: "your_user_name",
	Auth: []ssh.AuthMethod{
		ssh.Password("your_password")
	},
}
```

If you want to authenticate by using SSH certificate, there are two methods to obtain your ssh key:

SSH certificate file
You can parse your private key file by using ssh.ParsePrivateKey function. This is required by ssh.PublicKeys auth method function that creates a ssh.AuthMethod instance from private key.

```golang
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
```
Then you should instanciate ssh.ClientConfig:

```golang
sshConfig := &ssh.ClientConfig{
	User: "your_user_name",
	Auth: []ssh.AuthMethod{
		PublicKeyFile("/path/to/your/pub/certificate/key")	
	},
}
```

SSH agent
SSH Agent is a program that runs during user session in *nix system. It stores the private keys in an encrypted form. Because typing the passphrase can be tedious, many users would prefer to using it to store their private keys.

You can obtain all stored keys via SSH_AUTH_SOCK environment variable which stores the SSH agent unix socket. We should access the keys by calling net.Dial and then instanciate an agent client used by ssh.PublicKeysCallback factory auth method.

```golang
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
```

Then you can use the function to instanciate the client config in the following way:

```golang
sshConfig := &ssh.ClientConfig{
	User: "your_user_name",
	Auth: []ssh.AuthMethod{
		SSHAgent()
	},
}
```

Note that you can add your certificate to the SSH agent by using the following command:

$ ssh-add /path/to/your/private/certificate/file
Establishing new SSH connection
After we instaciated the ssh.ClientConfig object. We should be able to establish a new connection to the remote server by calling ssh.Dial:

```golang
connection, err := ssh.Dial("tcp", "host:port", sshConfig)
if err != nil {
	return nil, fmt.Errorf("Failed to dial: %s", err)
}
```

Creating a new session
After we established the connection, we should be able to open a new session that acts as an entry point to the remote terminal. We should use the connection in the following manner:

```golang
session, err := connection.NewSession()
if err != nil {
	return nil, fmt.Errorf("Failed to create session: %s", err)
}
```

Before we will be able to run the command on the remote machine, we should create a pseudo terminal on the remote machine. A pseudoterminal (or “pty”) is a pair of virtual character devices that provide a bidirectional communication channel.

We should create an xterm terminal that has 80 columns and 40 rows.

```golang
modes := ssh.TerminalModes{
	ssh.ECHO:          0,     // disable echoing
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
	session.Close()
	return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
}
```

If we want to attach our os.Stdin, os.Stdout and os.Stderr to remote command we should open pipes between the local process and remote process. Forthunatelly, ssh.Session object provides that out of the box by invoking session.Stdinpipe(), session.Stdoutpipe() and session.Stdouterr() functions. Then we should asyncronously copy the end of the pipes to the right file descriptors by using go routines and os.Copy function.

```golang
stdin, err := session.StdinPipe()
if err != nil {
	return fmt.Errorf("Unable to setup stdin for session: %v", err)
}
go io.Copy(stdin, os.Stdin)

stdout, err := session.StdoutPipe()
if err != nil {
	return fmt.Errorf("Unable to setup stdout for session: %v", err)
}
go io.Copy(os.Stdout, stdout)

stderr, err := session.StderrPipe()
if err != nil {
	return fmt.Errorf("Unable to setup stderr for session: %v", err)
}
go io.Copy(os.Stderr, stderr)
```

Command execution
Then we can execute the command on the remote machine by using session.Run function:

```golang
err = session.Run("ls -l $LC_USR_DIR")
// If we want to transfer some environment variable to the remote machine, we should use session.Setenv function to do that.

if err := session.Setenv("LC_USR_DIR", "/usr"); err != nil {
	return err
}
 ```
[1. 原文链接](https://blog.ralch.com/tutorial/golang-ssh-connection/)  
[2. ssh-agent and ssh-add](http://blog.joncairns.com/2013/12/understanding-ssh-agent-and-ssh-add/)
