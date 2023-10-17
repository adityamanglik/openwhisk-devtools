sudo apt update
sudo apt -y dist-upgrade
git config --global user.name "AM"
git config --global user.email "am5523@columbia.edu"
sudo apt install -y openjdk-8-jdk
sudo apt install -y nodejs npm zip
sudo apt install -y python3-pip
sudo apt install -y jq
pip install matplotlib
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh ./get-docker.sh
# docker pull openwhisk/action-nodejs-v10:nightly