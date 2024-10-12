sudo apt update
sudo apt -y dist-upgrade
sudo apt install --reinstall linux-firmware
git config --global user.name "AM"
git config --global user.email "aditya.account@outlook.com"
sudo sysctl net.ipv4.ip_local_port_range="15000 61000"
sudo sysctl net.ipv4.tcp_fin_timeout=30
sudo sysctl net.ipv4.tcp_tw_reuse=1
sudo sysctl -w fs.file-max=262144
sudo apt install python3-locust -y
sudo apt install apache2-utils  -y
sudo apt-get install gcc-multilib jq -y #for go
curl -LO get.golang.org/$(uname)/go_installer && chmod +x go_installer && ./go_installer --version $(curl https://go.dev/dl/?mode=json | jq -r '.[0].version') && rm go_installer
source /users/am_CU/.bash_profile
go get -u gonum.org/v1/gonum/...
sudo apt install python3-pip -y
pip install matplotlib requests pandas
go env -w GO111MODULE=off
go get -u gonum.org/v1/gonum/...
ssh-keygen -t ed25519 -C "am5523@columbia.edu"
