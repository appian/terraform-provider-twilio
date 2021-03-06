[![Build Status](https://travis-ci.com/Preskton/terraform-provider-twilio.svg?branch=master)](https://travis-ci.com/Preskton/terraform-provider-twilio)

# Twilio Terraform Provider

The goal of this Terraform provider plugin is to make manging your Twilio account easier.

Current features:

- `twilio_subaccount`
  - Create
  - Update
  - Delete
- `twilio_workspace`
  - Create
  - Update
  - Delete
- `twilio_phoneNumber`
  - Create
  - Update
  - Delete
- `twilio_worker`
  - Create
  - Update
  - Delete
- `twilio_taskQueue`
  - Create
  - Update
  - Delete
- `twilio_taskChannel`
  - Create
  - Update
  - Delete
- `twilio_workflow`
  - Create
  - Update
  - Delete
- `twilio_application`
  - Create
  - Update
  - Delete

More coming soon.

## Build
Run:
```
make plugin
```
to build and move the plugin to `~/.terraform.d/plugins` which is where terraform will look for all 3rd party plugins

## Getting Started

1. Start a trial account at twilio.com (if you don't have one already). Use the Console Dashboard to take note of your Account SID (a long string starts with `AC` and looks like a GUID) and Auth Token (also a long GUID-like string, hidden under the `View` link).
2. Download the latest release of the provider and place in your `~/.terraform.d/plugins` directory.
3. Use the example below, replacing `account_sid` and `auth_token` with the appropriate values.
4. `terraform apply` Note: this will cost you REAL MONEY (or at the very least trial credits).

## Debugging

Adding the following lines to your bash profile will enable additional logging.
```
export TF_LOG=TRACE
export TF_LOG_PATH=./terraform.log
export DEBUG_HTTP_TRAFFIC=true
```
## Example

Note: running and applying the below could cost you REAL MONEY! Please use this tool wisely!

```hcl
provider "twilio" {
    account_sid = "<your account sid here>"
    auth_token = "<your auth token here>"
}

resource "twilio_subaccount" "woomy" {
    friendly_name = "Woomy Subaccount #1"
}

resource "twilio_application" "new_twiml_app" {
    friendly_name = "My new TwiML application"
}

resource "twilio_worker" "test_worker" {
    friendly_name = "Your Name"
    workspace_sid = "WSXXXXXXXXXXXXXX"
}

resource "twilio_taskQueue" "normal_support" {
    friendly_name = "Normal Support"
    workspace_sid = "WSXXXXXXXXXXXXXX"
}

resource "twilio_workflow" "test_workflow" {
    friendly_name = "Test Workflow"
    workspace_sid = "WSXXXXXXXXXXXXXX"
    configuration = <<EOF
{
  "task_routing":{
    "filters": [
      {
        "filter_friendly_name": "Prioritizing Filter",
        "expression": "1==1",
        "targets": [
          {
            "queue": "WQccc",
            "priority": 1,
            "timeout": 300
          },
          {
            "priority": 10
          }
        ]
      }
    ], 
    "default_filter": {
      "queue": "WQccc"
    }
  }
}
EOF
}

resource "twilio_phoneNumber" "test_phone_number" {
    friendly_name = "Test Phone Number"
    search = "310"
    country_code = "US"
}

resource "twilio_taskChannel" "custom" {
  friendly_name = "Custom Channel"
  workspace_sid = "WSXXXXXXXXXXXXXX"
  unique_name = "custom"
  channel_optimized_routing = true
}

resource "twilio_workspace" "workspace" {
  friendly_name = "My Twilio workspace"
}
```
