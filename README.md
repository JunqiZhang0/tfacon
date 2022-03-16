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
##### Via pip
```bash
pip install "git+https://github.com/JunqiZhang0/tfacon.git@develop#egg=tfacon&subdirectory=pip_package"
```
### Get Started

#### tfacon.yml
This is where you store all parameters like this
__You must use auth_token from a superadmin account to run tfacon, otherwise the validation will fail!__
```yaml
launch_id: "<your launch_id>"
project_name: "<your project name>"
auth_token: "xxxxx-xxxx-xxxxx-xxxx-xxxxx" 
platform_url: "https://reportportal-<your_domain>.com"
tfa_url: "https://tfa-<your_domain>:443/"
connector_type: "RPCon"
```
The default name is tfacon.yml, you can change this by edit this environment variable __TFACON_YAML_PATH__ 

#### tfacon.cfg
This is where you put all the config information for tfacon
```ini
[config]
retry_times=2
concurrency=true
add_attributes=true
```
The default name for this cfg file is tfacon.cfg, you can change this by edit this environment variable __TFACON_CONFIG_PATH__

#### list
```bash
❯ tfacon list -h
list all information constructed from tfacon.yml/enviroment variables/cli

Usage:
  tfacon list [flags]

Flags:
      --auth-token string       The AUTH_TOKEN of report portal
      --connector-type string   The type of connector you want to use (example: RPCon, PolarionCon, JiraCon) (default "RPCon")
  -h, --help                    help for list
      --launch-id string        The launch id of report portal
      --platform-url string     The url to the test platform (example: https://reportportal-<your_domain>.com) (default "default value for platform url")
      --project-name string     The project name of report portal
      --tfa-url string          The url to the TFA Classifier (default "default value for tfa url")
```

Output Example:
```bash
❯ tfacon list
--------------------------------------------------
tfacon  1.0.0
Copyright (C) 2021, Red Hat, Inc.
-------------------------------------------------


2021/08/06 03:21:33 Printing the constructed information
LaunchID:        968
ProjectName:     TFA_RP
AuthToken:       xxxx-xxxx-xxxxxxx-xxxxxx-xxxxxxxxx
RPURL:           https://reportportal-<your_domain>.com
Client:          &{<nil> <nil> <nil> 0s}
TFAURL:          https://tfa-<your_domain>.com
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
      --platform-url string     The url to the test platform(example: https://reportportal-<your_domain>.com) (default "default value for platform url")
      --project-name string     The project name of report portal
      --tfa-url string          The url to the TFA Classifier (default "default value for tfa url")
```

Example Output
```bash
❯ tfacon run --project-name "project_name" --launch-id 1006

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
__You must use auth_token from a superadmin account to run tfacon, otherwise the validation will fail!__
```bash
❯ tfacon validate -h
validate if the parameter is valid and if the urls are accesible

Usage:
  tfacon validate [flags]

Flags:
      --auth-token string       The AUTH_TOKEN of report portal
      --connector-type string   The type of connector you want to use(example: RPCon, PolarionCon, JiraCon) (default "RPCon")
  -h, --help                    help for validate
      --launch-id string        The launch id of report portal
      --platform-url string     The url to the test platform(example: https://reportportal-<your_domain>.com) (default "default value for platform url")
      --project-name string     The project name of report portal
      --tfa-url string          The url to the TFA Classifier (default "default value for tfa url")

Global Flags:
  -v, --verbose   You can add this tag to print more detailed info
```

Example Output
```bash
❯ tfacon validate --project-name "TFACON" --launch-id 231
LaunchID:        231
ProjectName:     TFACON
AuthToken:       xxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx
RPURL:           https://reportportal-<your_domain>.com
Client:          &{<nil> <nil> <nil> 0s}
TFAURL:          https://tfa.com/latest/model

Validation Passed!
```
You can also add -v to have more detailed information for validation error
#### init
```bash
❯ tfacon init -h                                     
init will create a sample workspace for tfacon

Usage:
  tfacon init [flags]

Flags:
  -h, --help   help for init

Global Flags:
  -v, --verbose   You can add this tag to print more detailed info
```


### Advanced Config
__You can set up advanced config in ./tfacon.cfg by default or you can define the location of the config file with TFACON_CONFIG_PATH environment variable__

Example
```ini
[config]
concurrency=True
add_attributes=true
```

#### Set Concurrency
__You can set this to True or False, if you set it to True tfacon will deal with the test items in an async non-blocking way(faster), you can disable it to have a more clear view, but you will have a slower run compared to setting it to True__

#### add_attributes
__You can enable this to add an extra attribute "AI Prediction" to all the test items, the value of this attribute will be the prediction extracted from TFA Classifier__


## Developer Guide
### Archietrcue
#### UML graph
![uml](docs/image/tfacon_uml.svg)

## Contributor Guide
### Branch name
__*develop*__ is the development branch
__*master*__ is the stable branch

### More Details
#### Release Information
#### Video Tutorial
#### How To Embed tfacon to CI Pipelines
