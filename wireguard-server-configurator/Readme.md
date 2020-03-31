# Wireguard Server Configurator
Is a go binary which should run as cron every minute (dependes on the number of actions)

# Install (Ubuntu/Debian)
Download the binary
```
mkdir /opt/go-wg
cd /opt/go-wg
wget https://kekskurse-public.s3.eu-central-1.amazonaws.com/go-wg/wireguard-server-configurator
chown root:root /opt/go-wg/wireguard-server-configurator
chmod u+x /opt/go-wg/wireguard-server-configurator
```

Create a config file and put it to /etc/go-wg/server.yaml

```
DBConnectionString: "root:example@tcp(127.0.0.1:4306)/wg"
listenPort: 51820
ipRange: "10.42.133.0/24"
serverCertificatePath: "/etc/go-wg/server"
```

Change the parameter.

Create the follwong cron by run 
```
sudo su
crontab -e
```
```
* * * * *     /opt/go-wg/wireguard-server-configurator
```