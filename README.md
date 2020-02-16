## Install NATS (rpi)

Download NATS and install it in `/opt`:
```
wget https://github.com/nats-io/nats-server/releases/download/vx.y.z/nats-server-vx.y.z-linux-arm6.zip
unzip nats-server-vx.y.z-linux-arm6.zip 
sudo mkdir -p /opt/nats-server/bin
sudo cp ~/nats-server-vx.y.z-linux-arm6/nats-server /opt/nats-server/bin/
```

Create a dir for config:
```
sudo mkdir /etc/opt/nats-server/
```

Create a new file `/etc/systemd/system/nats.service`:
```
[Unit]
Description=NATS Server
After=network.target

[Service]
ExecStart=/opt/nats-server/bin/nats-server
WorkingDirectory=/etc/opt/nats-server
StandardOutput=inherit
StandardError=inherit
Restart=always
#User=nats

[Install]
WantedBy=multi-user.target
```

Ideally also create a new user _nats_, so that the service is not run as _root_.

Start the new service:
```
sudo systemctl start nats
```

Make sure it runs:
```
sudo systemctl status nats
```

Enable it so that it is started on boot:
```
sudo systemctl enable nats
```

Other useful commands:
```
sudo systemctl stop nats
sudo systemctl restart nats
```

# Install sunset

```
cd cmd/sunset
go build
sudo mkdir -p /opt/sunset/bin
sudo mkdir /etc/opt/sunset
sudo cp sunset /opt/sunset/bin
```

Create `/etc/systemd/system/sunset.service`
```
[Unit]
Description=Sunset events
After=nats.service

[Service]
ExecStart=/opt/sunset/bin/sunset
WorkingDirectory=/etc/opt/sunset
StandardOutput=inherit
StandardError=inherit
Restart=always
#User=nats

[Install]
WantedBy=multi-user.target
```

Start and enable it so that it is started on boot:
```
sudo systemctl start sunset
sudo systemctl enable sunset
```

# Install lightcontroller

```
cd cmd/lightcontroller
go build
sudo mkdir -p /opt/lightcontroller/bin
sudo mkdir /etc/opt/lightcontroller
sudo cp lightcontroller /opt/lightcontroller/bin
```

Create `/etc/systemd/system/lightcontroller.service`
```
[Unit]
Description=Light controller
After=nats.service

[Service]
ExecStart=/opt/lightcontroller/bin/lightcontroller
WorkingDirectory=/etc/opt/lightcontroller
StandardOutput=inherit
StandardError=inherit
Restart=always
#User=nats

[Install]
WantedBy=multi-user.target
```

Start and enable it so that it is started on boot:
```
sudo systemctl start lightcontroller
sudo systemctl enable lightcontroller
```

# Install controller

```
cd cmd/controller
go build
sudo mkdir -p /opt/controller/bin
sudo mkdir /etc/opt/controller
sudo cp controller /opt/controller/bin
```

Create `/etc/systemd/system/controller.service`
```
[Unit]
Description=Controller
After=nats.service

[Service]
ExecStart=/opt/controller/bin/controller
WorkingDirectory=/etc/opt/controller
StandardOutput=inherit
StandardError=inherit
Restart=always
#User=nats

[Install]
WantedBy=multi-user.target
```

Start and enable it so that it is started on boot:
```
sudo systemctl start controller
sudo systemctl enable controller
```
