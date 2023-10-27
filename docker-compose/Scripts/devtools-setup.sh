sudo apt update
sudo apt -y dist-upgrade
git config --global user.name "AM"
git config --global user.email "am5523@columbia.edu"
sudo apt install -y openjdk-8-jdk
sudo apt install -y nodejs npm zip
sudo apt install -y python3-pip
sudo apt install -y jq
sudo apt install python3-locust
pip install matplotlib requests
pip3 install paramiko
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh ./get-docker.sh
# docker pull openwhisk/action-nodejs-v10:nightly