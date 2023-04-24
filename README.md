
# OpenShift Network Validator
OpenShift network validator tool can be used to validate the networking configuration of an OpenShift cluster before deploying it. 
OpenShift Network Validator can be used to validate the DNS configuration, URL access, and network connectivity of an OpenShift cluster. OpenShift Network Validator can be run using the command-line interface, and it can also be used as a library in a Go program. OpenShift Network Validator is available for Linux, macOS, and Windows operating systems.

To contribute to the development of OpenShift Network Validator, you need to have Go and Git installed on your system. You can clone the repository and run the development commands to build and test OpenShift Network Validator.

[![Build](https://github.com/tosin2013/openshift-network-validator/actions/workflows/build.yml/badge.svg)](https://github.com/tosin2013/openshift-network-validator/actions/workflows/build.yml)
[![Release](https://github.com/tosin2013/openshift-network-validator/actions/workflows/release.yml/badge.svg)](https://github.com/tosin2013/openshift-network-validator/actions/workflows/release.yml)

## Running on Linux
```bash 
curl -OL https://github.com/tosin2013/openshift-network-validator/releases/download/v0.0.1/openshift-network-validator-v0.0.1-linux-amd64.tar.gz
curl -OL https://raw.githubusercontent.com/tosin2013/openshift-network-validator/main/install_types/sample.yaml
tar -xvf openshift-network-validator-v0.0.1-linux-amd64.tar.gz
chmod +x openshift-network-validator-linux-amd64
sudo mv openshift-network-validator-linux-amd64 /usr/local/bin/openshift-network-validator
openshift-network-validator dns  --config  sample.yaml --cluster-name mycluster --base-domain example.com
openshift-network-validator url-access  --config sample.yaml
openshift-network-validator test-networking  --config  sample.yaml
```

## Running on MacOS
```bash 
curl -OL https://github.com/tosin2013/openshift-network-validator/releases/download/v0.0.1/openshift-network-validator-v0.0.1-darwin-amd64.tar.gz
curl -OL https://raw.githubusercontent.com/tosin2013/openshift-network-validator/main/install_types/sample.yaml
tar -xvf openshift-network-validator-v0.0.1-darwin-amd64.tar.gz
chmod +x openshift-network-validator-v0.0.1-darwin-amd64
sudo mv openshift-network-validator-v0.0.1-darwin-amd64 /usr/local/bin/openshift-network-validator
openshift-network-validator dns  --config  sample.yaml --cluster-name mycluster --base-domain example.com
openshift-network-validator url-access  --config sample.yaml
openshift-network-validator test-networking  --config  sample.yaml
```

## Running on Windows
```powershell 
Invoke-WebRequest -Uri "https://github.com/tosin2013/openshift-network-validator/releases/download/v0.0.1/openshift-network-validator-v0.0.1-windows-amd64.exe" -OutFile "openshift-network-validator-v0.0.1-windows-amd64.exe"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/tosin2013/openshift-network-validator/main/install_types/sample.yaml" -OutFile "sample.yaml"
.\openshift-network-validator-v0.0.1-windows-amd64.exe  dns  --config  sample.yaml --cluster-name mycluster --base-domain example.com
.\openshift-network-validator-v0.0.1-windows-amd64.exe  url-access  --config sample.yaml
.\openshift-network-validator-v0.0.1-windows-amd64.exe  test-networking   --config  sample.yaml
```

## Deveploer requirements
* [Go](https://gist.github.com/tosin2013/d4f4420231a96aed2116efb4d6b151a0)
* git

```
git clone https://github.com/tosin2013/openshift-network-validator.git
cd openshift-network-validator
```

## Development commands
``` 
go run main.go dns --config install_types/sample.yaml --cluster-name mycluster --base-domain example.com
go run main.go url-access --config install_types/sample.yaml
go run main.go test-networking --config install_types/sample.yaml
```
