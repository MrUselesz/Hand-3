# Hand-in 3
In a terminal open the server by executing go run server.go in the server file. The server is now running.
In a separate terminal open the client file and run go run client.go  to open a client, it will request a name, write in terminal desired name this will the name displayed for all actions taken by that client. 
Once a client has been given a name it can send messages by simply typing in the terminal, the messages can have a maximum length of 128, a message longer than this will not be send and instead inform client that the message is too long, and a new message can then be written. Once a client wishes to leave the server it can do so typing in /leave in the terminal, it can re-enter the terminal by typing a new name (or the same one).
