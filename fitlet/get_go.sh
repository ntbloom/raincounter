GO_VERSION="go1.18.linux-amd64"
wget https://go.dev/dl/$GO_VERSION.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf $GO_VERSION.tar.gz