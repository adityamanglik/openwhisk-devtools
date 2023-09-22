# Install docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh ./get-docker.sh

# Fix docker permissions
# sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
docker run hello-world