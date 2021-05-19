package main

type server struct {
	rooms map[string]*room
	commands chan command
}

func newSever() *server {
	return &server{
		rooms: make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run(){
	for cmd := range s.commands{
		switch cmd.id {
		case CMD_NICKNAME:
			s.nickname(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.roomsList(cmd.client, cmd.args)
		case CMD_MESSAGE:
			s.message(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}
func (s *server) newClient(conn net.Conn) {
	log.Printf("New client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn: conn,
		nickname: "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) nickname(c *client, args []string){
	c.nickname = args[1]
	c.message(fmt.Sprintf("Your name is set to %s", c.nickname))
}
func (s *server) join(c *client, args []string){
	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name: roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the chat", c.nickname))
	c.message(fmt.Sprintf("You have joined %s", r.name))
}
func (s *server) roomsList(c *client, args []string){
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.message(fmt.Sprintf("The available rooms are: %s"), strings.Join(rooms, ", "))
}
func (s *server) message(c *client, args []string){
	if c.room == nil {
		c.err(errors.New("You are not in a chat"))
		return
	}
	c.room.broadcast(c, c,nickname + ": " + strings.Join(args[1:len(args)], " ")
}
func (s *server) quit(c *client, args []string){
	log.Printf("Client has been disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.coon.Close()
}
func (s *server) quitCurrentRoom(c *client){
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the chat", c.nickname)
	}
}