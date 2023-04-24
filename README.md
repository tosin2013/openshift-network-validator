
# OpenShift network validator

[![Build](https://github.com/tosin2013/openshift-network-validator/actions/workflows/build.yml/badge.svg)](https://github.com/tosin2013/openshift-network-validator/actions/workflows/build.yml)

## Deveploer requirements - WIP
* [Go](https://gist.github.com/tosin2013/d4f4420231a96aed2116efb4d6b151a0)
* git

```
git clone https://github.com/tosin2013/openshift-network-validator.git
cd openshift-network-validator
```

### run app
``` 
go run main.go dns --config install_types/sample.yaml --cluster-name mycluster --base-domain example.com
go run main.go url-access --config install_types/sample.yaml
go run main.go test-networking --config install_types/sample.yaml
```
