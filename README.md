# `dodyndns`: Dynamic DNS updater for DigitalOcean

This is the simple tool I wrote to keep my DigitalOcean DNS records up to date
on my home machines. It was slapped together in about an hour, so it's got a
number of thigs to fix:

 * Root records for the domain have to be specified as @.my.domain.com, as
   the program currently naively splits addresses on dots. 
 * If the host has multiple valid public IPv4 or IPv6 addresses, the program
   currently only takes the first address it finds from each family.
 * It only makes HTTP requests for IPv4 public addresses behind NAT - Not sure
   I really care to fix this one, as I'm against v6 NAT in general, and IPv6
   temporary addresses futher complicate this. I'm happy to hear ideas about
   how this should work, though.

## Installation

1. Install [Go](http://golang.org)
2. Create a [workspace](https://golang.org/doc/code.html#Workspaces)
3. `go get github.com/justinian/dodyndns`
4. The tool is installed at `$GOPATH/bin/dodyndns`
