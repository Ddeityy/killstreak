#!/usr/bin/bash
echo "Killing old service"
systemctl --user -M $SUDO_USER@ stop killstreak.service
echo "Copying the executables"
cp bin/killstreak /home/$SUDO_USER/.local/share/
echo "Creating service file"
sudo touch /etc/systemd/user/killstreak.service
sudo bash -c 'cat' << EOF > /etc/systemd/user/killstreak.service --cut=$cut
[Unit]
Description=Killstreak service
[Service]
ExecStart=/home/$SUDO_USER/.local/share/killstreak
Restart=always
[Install]
WantedBy=default.target
EOF
echo "Reloading systemd daemon"
systemctl --user -M $SUDO_USER@ daemon-reload
echo "Enabling and starting the service"
echo ---------------------------------------------
echo | cat /etc/systemd/user/killstreak.service
echo ---------------------------------------------
systemctl --user -M $SUDO_USER@ enable --now killstreak.service
systemctl --user -M $SUDO_USER@ status killstreak.service