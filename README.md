wesweard
========

A very simple chat-like TCP server that can be used to broadcast messages across a group of listeners. There is no state, messages sent to the broadcast socket are relayed to any clients connected to the receive socket.

This can easily be run using docker:
```
$ docker build -t wesweard .
$ docker run -d -p 4444:4444 -p 5555:5555 wesweard bash -c '/wesweard/bin/wesweard --port-recv 4444 --port-bcast 5555'
```

To broadcast:
```
$ telnet localhost 4444
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
{your message}
^]
telnet> quit
Connection closed.
```

To listen:
```
$ telnet localhost 4444
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
{your message}
^]
telnet> quit
Connection closed.
```

development
===========

This project uses gb and all vendored dependencies are included.
