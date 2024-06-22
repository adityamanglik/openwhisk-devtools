sudo apt update
sudo apt -y dist-upgrade
sudo apt-get install --reinstall linux-firmware
git config --global user.name "AM"
git config --global user.email "am5523@columbia.edu"
sudo apt install -y openjdk-8-jdk
sudo apt install -y nodejs npm zip
sudo apt install -y python3-pip
sudo sudo sysctl -w fs.file-max=262144
sudo sysctl net.core.somaxconn=1024
sudo sysctl net.core.netdev_max_backlog=2000
sudo sysctl net.ipv4.tcp_max_syn_backlog=2048
sudo apt-get install gcc-multilib #for go
wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
sudo tar -xvf go1.22.4.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
source ~/.profile
pip install pillow
pip install psutil # Python server
# Extend file handle limits
echo '* soft nofile 10000' | sudo tee -a /etc/security/limits.conf
echo '* hard nofile 10000' | sudo tee -a /etc/security/limits.conf
echo 'session required pam_limits.so' | sudo tee -a /etc/pam.d/common-session
echo 'session required pam_limits.so' | sudo tee -a /etc/pam.d/common-session-noninteractive
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh ./get-docker.sh
sudo apt install -y docker-compose
sudo reboot now