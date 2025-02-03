<p align="center">
  <img width="300" height="300" src="images/logo.webp">
</p>

<p align="center" style="margin: 0; padding: 0;">
  <a href="https://pkg.go.dev/github.com/mtnmunuklu/bridge">
    <img src="https://img.shields.io/badge/%F0%9F%93%9A%20godoc-pkg-informational.svg" alt="Go Doc">
  </a> 
  <a href="https://goreportcard.com/report/github.com/mtnmunuklu/bridge">
    <img src="https://img.shields.io/badge/%F0%9F%93%9D%20goreport-A+-success.svg" alt="Go Report">
  </a> 
  <a href="https://travis-ci.com/">
    <img src="https://img.shields.io/badge/%E2%9A%99%20build-X-success.svg" alt="Build Status">
  </a>
</p>

# Bridge

Bridge is a versatile tool that connects **Sigma rules** to **SPL (Splunk Processing Language)**, enabling seamless integration between the two for more efficient security analysis and threat detection.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
  - [Normal Installation](#normal-installation)
  - [Docker Installation](#docker-installation)
- [Usage](#usage)
  - [Normal Usage](#normal-usage)
  - [Docker Usage](#docker-usage)
- [Demo](#demo)
- [Contributing](#contributing)
- [Acknowledgement](#acknowledgement)
- [License](#license)

## Overview

Bridge is designed to simplify the process of converting **Sigma** rules into **SPL** queries for use in **Splunk** environments. Sigma is an open-source rule format for creating and sharing detection rules, while SPL is the powerful query language used by Splunk to process and analyze log data. With Bridge, security analysts can easily convert Sigma-based detection rules into Splunk queries, improving their detection capabilities without reinventing the wheel.

Bridge supports **Sigma to SPL** conversion and ensures you get the most out of both Sigma’s rule sets and Splunk's querying power.

## Installation

To use Bridge, you can install it in two ways:

### Normal Installation

Bridge offers precompiled ZIP files for multiple platforms. Download the appropriate ZIP file for your platform:

- [Windows](https://github.com/yourusername/bridge/releases/latest/download/bridge-windows-latest.zip)
- [Linux](https://github.com/yourusername/bridge/releases/latest/download/bridge-linux-latest.zip)
- [macOS](https://github.com/yourusername/bridge/releases/latest/download/bridge-macos-latest.zip)

After downloading, extract the ZIP file to a folder of your choice. Make sure the directory containing the Bridge executable is included in your system’s PATH so you can run it from the command line.

### Docker Installation

For a containerized solution, you can run Bridge using **Docker**. Docker simplifies the setup and ensures a consistent environment.

To get started with Docker, ensure Docker is installed on your machine. If it's not installed yet, follow the instructions on the official site: [https://www.docker.com/get-started](https://www.docker.com/get-started).

Once Docker is ready, follow these steps to set up Bridge:

1. **Clone the Repository**:

   ```shell
   git clone https://github.com/mtnmunuklu/bridge.git
   ```
2. **Navigate to Docker Directory**: Go to the docker directory inside the cloned repository:

   ```shell
   cd tools/docker
   ```
3. **Build Docker Image and Start Container**: Use the setup script to build the Docker image named bridge-image:

   ```shell
   go run setup_docker_bridge.go -rules <rulesDirectory> -config <configFile> -output <outputDirectory>
   ```
   
   This script will handle the building of the Docker image and starting the container for you.

That's it! You have successfully installed Bridge on your system. You can now proceed to the [Usage](#usage) section to learn how to use Bridge.

If you prefer to build Bridge from source, you can refer to the [Build Instructions](BUILD.md) for detailed steps on how to build and install it on your platform.

## Usage

To use Bridge, you’ll need Sigma rules in YAML format. Sigma rules can be found in the Sigma GitHub repository: https://github.com/Neo23x0/sigma/tree/master/rules. The configuration file for Splunk should be obtained from your system administrator.

### Normal Usage

To convert Sigma rules to SPL in a local environment, use the following command:

1. Prepare Sigma rules and a configuration file.
2. Convert Sigma rules to the query language of CRYPTTECH's SIEM product by running the following command:

    ```shell
    ./bridge -filepath <path-to-sigma-rules> -config <path-to-config> [-json] [-output <output-directory>]
    ```
    or
    ```shell
    ./bridge -filecontent <content-to-sigma-rules> -configcontent <content-to-config> [-json] [-output <output-directory>]

    ```

### Docker Usage

If you have installed Bridge using Docker, you can use the following command to run Bridge inside the Docker container:

```shell
docker exec bridge ./bridge -filepath <path-to-sigma-rules> -config <path-to-config> [-json] [-output <output-directory>]
```
or
```shell
docker exec bridge ./bridge -filecontent <content-to-sigma-rules> -configcontent <content-to-config> [-json] [-output <output-directory>]
```

The `filepath` flag specifies the location of the Sigma rules. This can be a file or directory path.

The `filecontent` flag allows you to provide the Base64-encoded content of Sigma rules directly as a string.

The `config` flag specifies the location of the configuration file for SPLUNK product.

The `configcontent` flag allows you to provide the Base64-encoded content of the configuration file directly as a string.

The `json` flag indicates that the output should be in JSON format.

The `output` flag specifies the directory where the output files should be written.

If the `json` flag is provided, Bridge will convert the Sigma rules to JSON format. If the `output` flag is provided, Bridge will save the output files to the specified directory. If neither flag is provided, the output will be displayed in the console.

## Contributing

Contributions to Bridge are welcome and encouraged! Please read the [contribution guidelines](CONTRIBUTING.md) before making any contributions to the project.

## Acknowledgements

We want to express our gratitude to the creators of the Sigma as their rule formats form the backbone of this project. More information about Sigma can be found [here](https://github.com/Neo23x0/sigma).

## License

Bridge is licensed under the MIT License. See [LICENSE](LICENSE) for the full text of the license.
