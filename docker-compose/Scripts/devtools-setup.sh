sudo apt update
sudo apt -y dist-upgrade
sudo apt install -y openjdk-8-jdk
sudo apt install -y nodejs npm zip
sudo apt install -y docker-compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh ./get-docker.sh
docker pull openwhisk/action-nodejs-v10:nightly