#Async_Terraform

### Introduction

At some point I had to run a lot o terraform one-by-one in a mono repo to manage simple resources in AWS. To address this, I did create this action.

### Overview

The ideia is to specify the project names containing `.tf` files under the `task` string while using the async_terraform action. (there is an example below).

Then the action will run all projects concurrently based on the number os workers configured.

### Usage

```
The Action must have three parameters
1. workers -> Specifies the number of workers running concurrently in a worker pool to execute the terraform tasks. The default value is 2.
2. action_name -> Specifies the name of the .yaml/.yml file located under `.github/workflows/`
3. verb -> Specifies the action to be performed by Terraform. Plan, apply or destroy. This is set manually by input using workflow_dispatch as the example below
4. tasks -> The list of tasks (terraform projects) to be read by the action
```

**Below is the folder structure example to use the action**
Remember to replace each string in "tasks" with each name of your project. In the example below I'm using terraform_1, terraform_2, terraform_3, terraform_4

```
├── .github
│   └── workflows
│       └── action.yml
├── terraform_1
│   ├── providers.tf
│   └── vpc.tf
├── terraform_2
│   ├── providers.tf
│   └── vpc.tf
├── terraform_3
│   ├── providers.tf
│   └── vpc.tf
├── terraform_4
│   ├── providers.tf
│   └── vpc.tf

```

### The action need to be like this

**The credentials step below can be modified to autenticate in another cloud providers**

```yaml
on:
  workflow_dispatch:
    inputs:
      verb:
        description: "Plan, apply or destroy"
        required: true
        type: choice
        options:
          - plan
          - apply
          - -----
          - destroy

name: Async_terraform
jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: AWS Credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Terraform Task
        uses: adalbertjnr/async_terraform@v1
        with:
          workers: 2
          action_name: action.yml
          verb: ${{ inputs.verb }}
          tasks: |
            terraform_1,
            terraform_2,
            terraform_3,
            terraform_4
```
