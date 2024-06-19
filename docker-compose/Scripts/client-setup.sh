sudo apt update
sudo apt -y dist-upgrade
sudo apt install --reinstall linux-firmware
git config --global user.name "AM"
git config --global user.email "am5523@columbia.edu"
sudo apt install python3-locust -y
sudo apt install apache2-utils  -y
sudo apt-get install gcc-multilib #for go
wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
sudo tar -xvf go1.22.4.linux-amd64.tar.gz
sudo mv go /usr/local
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.bashrc
echo "export GO111MODULE=on" >> ~/.bashrc
source ~/.profile
go get -u gonum.org/v1/gonum/...
sudo apt install python3-pip -y
pip install matplotlib requests pandas
go env -w GO111MODULE=off
go get -u gonum.org/v1/gonum/...
ssh-keygen -t ed25519 -C "am5523@columbia.edu"
