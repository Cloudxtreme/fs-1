package console

import (
	"fmt"
	"sync"
)

type Cmd struct {
	argv0 string
	cmd   func(*Cons, []string) error
}

var clibox struct {
	lock sync.Mutex
	cmd  []Cmd
}

func Exec(cons *Cons, buf string) error {
	argv := tokenize(buf)

	if len(argv) == 0 || argv[0][0] == '#' {
		return nil
	}

	clibox.lock.Lock()
	for _, c := range clibox.cmd {
		if c.argv0 == argv[0] {
			clibox.lock.Unlock()
			if err := c.cmd(cons, argv); err != nil /* && err != EUsage */ {
				return err
			}
			return nil
		}
	}
	clibox.lock.Unlock()

	return fmt.Errorf("%s: - eh?", argv[0])
}

func AddCmd(argv0 string, cmd func(*Cons, []string) error) error {
	clibox.lock.Lock()
	defer clibox.lock.Unlock()

	for _, c := range clibox.cmd {
		if argv0 == c.argv0 {
			return fmt.Errorf("cmd %q already registered", c.argv0)
		}
	}

	c := Cmd{
		argv0: argv0,
		cmd:   cmd,
	}
	clibox.cmd = append(clibox.cmd, c)
	return nil
}
