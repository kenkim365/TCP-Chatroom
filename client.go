package main

import "net"

type client struct {
	conn net.Conn
	nickname string
	room *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		message, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		message = strings.Trim(message, "\r\n")

		args := strings.Split(message, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
			case "/nickname":
				c.commands <- command{
					id: CMD_NICKNAME,
					clients: c,
					args: args,
				}
			case "/join":
				c.commands <- command{
					id: CMD_JOIN,
					clients: c,
					args: args,
				}
			case "/rooms":
				c.commands <- command{
					id: CMD_ROOMS,
					clients: c,
					args: args,
				}
			case "/message":
				c.commands <- command{
					id: CMD_MESSAGE,
					clients: c,
					args: args,
				}
			case "/quit":
				c.commands <- command{
					id: CMD_QUIT,
					clients: c,
					args: args,
				}
			default:
				c.err(fmt.Errorf("Unknown command: %s", cmd))

		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

func (c *client) message(message string) {
	c.conn.Write([]byte("> " + message + "\n"))
}