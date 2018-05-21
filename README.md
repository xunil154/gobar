# [GOBar](http://lmgtfy.com/?q=gobar)

Learning GO with shellz

**WARNING:** A lot of this code might look like actual *gobar*, but the intent
is to also be usefull for playing with other people's *gobar* systems.

# Structure

There are three main components to GOBar, 

* Core
* Interface
* Agent

```
        +-------+         +-------+
        | Agent +---+ +---+ Agent |
        +-------+   | |   +-------+
                    | |
                    | |
                    | |
               +----v-v---+
               |          |
       +------->   Core   <-------+
       |       |          |       |
       |       +-----^----+       |
       |             |            |
       |             |            |
+------+--+    +-----+---+   +----+----+
|Interface|    |Interface|   |Interface|
+---------+    +---------+   +---------+

```

## Core

A core service listens on 1337 that *Interfaces* and *Agents* interact with.

Multiple *Interfaces* can connect and disconnect as necessary (think like tmux)

The *Core* can then be configured by an *Interface* to listen for incoming 
connections and will monitor status until you decide to interact with them
from an *Interface*.

It will maintain shells **AND** log the sessions, so you will always keep a
record of what happened

The *Core* will also provide several additional services, such as 

* HTTP File Server to serve payloads
* Launch attacks

## Interface

The *Interface* component is what you actually interact with. It will connect
to the *Core* to send commands, receive data, and interact with connected
shells

## Agent

The *Agent* is a binary deployed to a server that connects back to the *Core*.

An *Agent* can be as simple as a reverse bash shell.

Eventually I hope to make these more complex

# Background

I am going through the [OSCP](https://www.offensive-security.com/information-security-certifications/oscp-offensive-security-certified-professional/)
certification, and they don't let you use a lot of the standard tools out there
to force you to learn from the ground up (which I appriciate). So being me, I'm
not going to go the easy route of a simple `nc -lkv 4444` to catch reverse
shells, because too many times I've accidentally hit `Ctrl+C` after working
hours trying to get the shell in the first place.

So for the OSCP, I set two additional goals for myself

1) Learn GO (thus this project)
2) Force myself to learn Radare2
