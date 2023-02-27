# GoSniff

Meant to listen messages in discord channels and forward them to another.



# General Flow
A single listener can listen to X number of channels across multiple servers. Each listener
opens up a connection to discord and then has to filter the messages only to what we care about.
Every new message fires a "message create" event that we must process. If we care about it, we clean
it up, transform the names and make it look nice and then send it off to the forwarder to relay it.
Forwarders have two modes currently, chatting as a user, and chatting as a bot using webhooks. Chatting
as a user requires less setup, but grants us far less control over the message. Using webhooks allow us 
to imitate it coming from the user that posted it. 