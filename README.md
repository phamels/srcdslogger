# SRCDSLOGGER
A simple logparser for Source Dedicated Server (Mainly Counter-Strike: Source at this point) written in Go.

It starts a UDP Server (configurable in `config.json`) that can receive logs from the Source Dedicated Server (e.g with `logaddress_add x.x.x.x:27500`)
Parses the logs and writes the stats to a MySQL database (SQL migration included).

Has a basic point system in place and has a lot of room for improvement :-)

Feel free to contriubute :-)

