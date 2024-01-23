sudo apt update
sudo apt -y dist-upgrade
sudo apt install --reinstall linux-firmware
git config --global user.name "AM"
git config --global user.email "am5523@columbia.edu"
sudo apt install python3-locust -y
sudo apt install jq -y
sudo apt install golang-go -y
sudo apt install python3-pip -y
sudo apt install apache2-utils  -y
pip install matplotlib requests
go env -w GO111MODULE=off
go get -u gonum.org/v1/gonum/...
ssh-keygen -t ed25519 -C "am5523@columbia.edu"
