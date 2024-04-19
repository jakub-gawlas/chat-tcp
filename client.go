package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

// Client connected to chat-tcp
type Client struct {
	conn net.Conn
	id   string
}

func NewClient(conn net.Conn) Client {
	return Client{
		conn: conn,
		id:   conn.RemoteAddr().String(),
	}
}

// HandleConnection read messages send by client, should be emitted to channel 'messages'
func (c *Client) HandleConnection(messages chan<- Message) {
	defer c.Close()

	// [1st shortcoming] I'm aware that character is not the same as byte (one char can contain more than one byte).
	// [2nd shortcoming] I assumed that any received message ends new line character,
	// 					 of course new line character can be in the middle of message,
	//					 then original message probably should be split to many separate messages.
	buffer := make([]byte, 128)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			// todo: handle other cases when client is disconnected and should stop process the one
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			c.Log("error while read message: " + err.Error())
			continue
		}

		data := buffer[:n]
		messages <- Message{
			data:     data,
			senderId: c.id,
		}
	}
}

// Send data to client
func (c *Client) Send(data []byte) error {
	// todo: set deadline to write message, some cleanup should be implemented,
	//		 client should be disconnected in case of eg. timeout
	if _, err := c.conn.Write(data); err != nil {
		return err
	}
	return nil
}

// Log message related to client
func (c *Client) Log(msg string) {
	fmt.Printf("[%s]: %s\n", c.id, msg)
}

// Close client connection
// todo: at now client.Close is invoked many times while terminating server, it generates unnecessary logs
func (c *Client) Close() {
	c.Log("client disconnected")
	if err := c.conn.Close(); err != nil {
		c.Log("error while close client: " + err.Error())
	}
}
