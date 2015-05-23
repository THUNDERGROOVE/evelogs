# evelogs

evelogs is a command line utiltiy to query through your EVE online log files.

I wrote it because I grew tired of manually going through my log folder to look for something that was said months ago

# Usage

evelogs as a couple of flags

`-l` Lists information about your chat logs as a whole;  Prints how many log files you have the count of lines of chat and how many characters have lines of chat in your logs

`-c` Channel can be used to query from a specific channel

`-u` User can be used to filter out messages not written by the specified user

`-r` Regex/filter if given a regex will filter out all messages that don't match it.  If given just a string we test if the message contains the given string
