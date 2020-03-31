# GO-WG

Some Software configure a Wireguard "Server" that Clients can easy connect to it. 
The Clients connect to the Server via a HTTP-API, send a Public Key and ask for access. 
If the Administrator approve the Access (via a WebGUI) the client get a IP-Address which it can use to connect to the Wireguard Server.

![Plan](statik/Plan.png "Plan")

* wg -> Package to config Device and Wireguard
* wireguard-server-configurator -> Cron runs on the Wireguard Server
* wireguard-server-http -> HTTP Interface runs on the Wireguard Server
* wireguard-client -> Client run as Cron on the Client which should connect to the Server

# Installataion
1) Install MariaDB and execute all scripts from the [database Folder](https://github.com/kekskurse/go-wg/tree/master/database) 
2) Install the [server configurator cron](https://github.com/kekskurse/go-wg/tree/master/wireguard-server-configurator)
3) Install the http server
4) Install a Client on another Server