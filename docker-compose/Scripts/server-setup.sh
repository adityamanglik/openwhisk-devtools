sudo apt update
sudo apt -y dist-upgrade
sudo apt-get install --reinstall linux-firmware
git config --global user.name "AM"
git config --global user.email "aditya.account@outlook.com"
sudo apt install -y openjdk-8-jdk
sudo apt install -y nodejs npm zip
sudo apt install -y python3-pip
sudo sysctl -w fs.file-max=262144
sudo sysctl net.core.somaxconn=1024
sudo sysctl net.core.netdev_max_backlog=2000
sudo sysctl net.ipv4.tcp_max_syn_backlog=2048
sudo apt-get install gcc-multilib jq -y #for go
curl -LO get.golang.org/$(uname)/go_installer && chmod +x go_installer && ./go_installer --version $(curl https://go.dev/dl/?mode=json | jq -r '.[0].version') && rm go_installer
source /users/am_CU/.bash_profile
pip install pillow
pip install psutil # Python server
# Extend file handle limits
echo '* soft nofile 10000' | sudo tee -a /etc/security/limits.conf
echo '* hard nofile 10000' | sudo tee -a /etc/security/limits.conf
echo 'session required pam_limits.so' | sudo tee -a /etc/pam.d/common-session
echo 'session required pam_limits.so' | sudo tee -a /etc/pam.d/common-session-noninteractive
# curl -fsSL https://get.docker.com -o get-docker.sh
# sudo sh ./get-docker.sh
# sudo apt install -y docker-compose
sudo reboot now
