# Terraform Google AppScript Provider

This is a [Terraform][terraform] provider for managing Google appscript projects.


## Installation

1. Download the latest compiled binary from [GitHub releases][releases].

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins/linux_amd64
    $ mv terraform-provider-googleappscript $HOME/.terraform.d/plugins/terraform-provider-googleappscript_v0.1.0
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```

    This will find the plugin locally.

## Usage

1. Since google appscript api does not support service account authentication,
you will need a OAuth token file to run provider
and unfortunately you will have to renew your token file whenever its expired.

    1. Go to APIs & Services of Google Cloud Platform
    1. Create a credential and download as a file
    1. Run this [script][script] with credentials file
    1. Copy-paste url to your browser
    1. Give permissions and copy back code to terminal
    1. Your token file should be created.

1. Create a Terraform configuration file:

    ```hcl
       provider "googleappscript" {
         token_file = "token.json"
       }

       resource "googleappscript_project" "example" {
         title = "terraform-example"

         script {
           name = "appsscript"
           type = "JSON"
           source = "{\"timeZone\":\"America/New_York\",\"exceptionLogging\":\"CLOUD\"}"
         }

         script {
           name = "hello"
           type = "SERVER_JS"
           source = "function helloWorld() {\n  console.log('goodbye, world!');}"
         }

       }
    ```

1. Run `terraform init` to pull in the provider:

    ```sh
    $ terraform init
    ```

1. Run `terraform plan` and `terraform apply` to create events:

    ```sh
    $ terraform plan

    $ terraform apply
    ```

[terraform]: https://www.terraform.io/
[releases]: https://github.com/pasali/terraform-provider-googleappscript/releases
[script]: https://gist.github.com/pasali/09cebb2599dfde5a5eef24b8e805b434
