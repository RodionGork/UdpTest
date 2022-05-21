# Go UDP Test

Just to recollect some basics on UDP, golang etc. Here is a basic app which
receives some "events", stores them and allows to view them somehow.

Reception is going to work via UDP.

Storing is 2-tier, "short memory" - say, for about 10 last minutes - as a
RAM structure - while in parallel events should be stored persistently
(I'm curious to try what speed SQLite can handle).

Viewing/visualizing to be provided by browsing storage via web-interface.

### Initial

With fist commit we get a test prototype which shows, for example, that
1000 UDP messages could be easily ingested per 1 second (on typical
contemporary hardware).