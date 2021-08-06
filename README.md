# Test Failure Classifier Connector 

## Description

## User Guide
### Installation
##### Via go get
```bash
go get -u github.com/JunqiZhang0/tfacon
```
##### Via go get(To be added)
```bash
pip install tfacon
```
### Get Started
#### list
#### run
#### validate
#### init(To be added)
### Advanced config
__You can set up advanced config in ./tfacon.cfg by default or you can define the location of the config file with TFACON_CONFIG_PATH environment variable__

Example
```ini
[config]
retry_times=2
concurrency=True
```

#### Set Concurrency
__You can set this to True or False, if you set it to True tfacon will deal with the test items in an async non-blocking way(faster), you can disable it to have a more clear view, but you will have a slower run compared to setting it to True__
#### Set retry(To be added)

## Developer Guide
### Archietrcue
#### UML graph
![uml](doc/image/tfacon_uml.svg)
#### OOD brief explanation
### concurrency

## Contributor Guide
### Branch name

