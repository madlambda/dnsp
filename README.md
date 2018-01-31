# DNSP

DNSP stands for DNS Poisoning. Actually what it does is
**MUCH** simpler than actual DNS poisoning, it will make
local resolutions on your host for a specific domain to
resolve to localhost, running at the same time an http
server on the localhost and providing static content
based on a provided directory.

It is very simple, calling it DNS poisoning is more
to make me feel fancy...I like to feel fancy =).

When you access a site by a domain name it will actually
be answered by dnsp running on your host. For use cases
use your imagination, it is barely useful.

You will need to run it as sudo since it will make changes
to /etc/hosts and will listen on the 80 port.
