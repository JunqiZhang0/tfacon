# Test Failure Classifier Connector 

## Description
__tfacon__ is a CLI tool for connecting Test Management Platforms and __Test Failure Analysis Classifier__. __Test Failure Analysis Classifier__ is
an AI/ML predictioner developed by *Red Hat D&O Data Science Team* which can predict the test's catalog. It supports __AutomationBug, ProductBug, SystemBug__ on Report Portal now. tfacon only support report portal at this moment. We will support more platforms in the future.
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
```bash
❯ tfacon list -h
list all information constructed from tfacon.yml/enviroment variables/cli

Usage:
  tfacon list [flags]

Flags:
      --auth-token string       The AUTH_TOKEN of report portal
      --connector-type string   The type of connector you want to use(example: RPCon, PolarionCon, JiraCon) (default "RPCon")
  -h, --help                    help for list
      --launch-id string        The launch id of report portal
      --platform-url string     The url to the test platform(example: https://reportportal-ccit.apps.ocp4.prod.psi.redhat.com) (default "default val for platform url")
      --project-name string     The project name of report portal
      --tfa-url string          The url to the TFA Classifier (default "default val for tfa url")
```

Output Example:
```bash
--------------------------------------------------
tfacon  0.0.1
Copyright (C) 2021, Red Hat, Inc.
-------------------------------------------------


2021/08/06 03:21:33 Printing the constructed information
LaunchId:        968
ProjectName:     TFA_RP
AuthToken:       xxxx-xxxx-xxxxxxx-xxxxxx-xxxxxxxxx
RPURL:           https://reportportal.com
Client:          &{<nil> <nil> <nil> 0s}
TFAURL:          https://tfa.com
```
#### run
```bash
❯ tfacon run -h                                       
run the info retrival and get the pridiction from TFA

Usage:
  tfacon run [flags]

Flags:
      --auth-token string       The AUTH_TOKEN of report portal
      --connector-type string   The type of connector you want to use(example: RPCon, PolarionCon, JiraCon) (default "RPCon")
  -h, --help                    help for run
      --launch-id string        The launch id of report portal
      --platform-url string     The url to the test platform(example: https://reportportal-ccit.apps.ocp4.prod.psi.redhat.com) (default "default val for platform url")
      --project-name string     The project name of report portal
      --tfa-url string          The url to the TFA Classifier (default "default val for tfa url")
```

Eample Output
```bash
❯ ./tfacon run --project-name "project_name" --launch-id 1006

2021/08/06 03:46:59 Getting prediction of test item(id): 54799
2021/08/06 03:46:59 Getting prediction of test item(id): 54900
2021/08/06 03:46:59 Getting prediction of test item(id): 54106
2021/08/06 03:46:59 Getting prediction of test item(id): 54555
2021/08/06 03:46:59 Getting prediction of test item(id): 54986
2021/08/06 03:46:59 Getting prediction of test item(id): 54411
2021/08/06 03:46:59 Getting prediction of test item(id): 54824
2021/08/06 03:46:59 Getting prediction of test item(id): 54642
2021/08/06 03:46:59 Getting prediction of test item(id): 54841
This is the return info from update: [{"issueType":"ab001","comment":"","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"ab001","comment":"","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"ab001","comment":"Should be marked with custom defect type","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"si001","comment":"","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"ab001","comment":"Should be marked with custom defect type","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"ab001","comment":"Should be marked with custom defect type","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"ab001","comment":"Should be marked with custom defect type","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"ab001","comment":"Should be marked with custom defect type","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]},{"issueType":"pb001","comment":"Should be marked with custom defect type","autoAnalyzed":false,"ignoreAnalyzer":false,"externalSystemIssues":[]}]
```
#### validate
#### init(To be added)
### Advanced Config
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

