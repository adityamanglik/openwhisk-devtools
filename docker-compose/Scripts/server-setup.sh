sudo apt update
sudo apt -y dist-upgrade
sudo apt-get install --reinstall linux-firmware
git config --global user.name "AM"
git config --global user.email "am5523@columbia.edu"
sudo apt install -y openjdk-8-jdk
sudo apt install -y nodejs npm zip
sudo apt install -y python3-pip
sudo apt install -y golang-go
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