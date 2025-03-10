# SugarFree

Less sugar (entropy) for your binaries

<p align="center">
  <img width="450" height="450" src="/Pictures/Logo.png"><br /><br />
  <!--<img alt="GitHub License" src="https://img.shields.io/github/license/nickvourd/SugarFree?style=social&logo=GitHub&logoColor=purple">
  <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/nickvourd/SugarFree?logoColor=yellow"><br />
  <img alt="GitHub forks" src="https://img.shields.io/github/forks/nickvourd/SugarFree?logoColor=red">
  <img alt="GitHub watchers" src="https://img.shields.io/github/watchers/nickvourd/SugarFree?logoColor=blue">
  <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/nickvourd/SugarFree?style=social&logo=GitHub&logoColor=green">-->
</p>

## Description

SugarFree is an open-source tool designed to analyze and reduce the entropy of a provided PE file.

![Static Badge](https://img.shields.io/badge/Go-lang-cyan?style=flat&logoSize=auto)
![Static Badge](https://img.shields.io/badge/Version-2.0%20(Ocean%20Words)-red?link=https%3A%2F%2Fgithub.com%2Fnickvourd%2FSugarFree%2Freleases)

SugarFree uses two different techniques (strategies) to reduce the entropy of a PE file:

- **zero**: Appends null bytes (`0x00`) to the end of the file, increasing its size while lowering entropy.  
- **word**: Appends random English words in byte format to the end of the file, increasing its size while lowering entropy.

The following list explains the meaning of each SugarFree command:

- **info**: Calculates the entropy of a PE file and its sections.
- **free**: Lowers the overall entropy of a PE file.

SugarFree is written in Golang, a cross-platform language, enabling its use on both Windows and Linux systems.

> If you find any bugs, don’t hesitate to [report them](https://github.com/nickvourd/SugarFree/issues). Your feedback is valuable in improving the quality of this project!

## Disclaimer

The authors and contributors of this project are not liable for any illegal use of the tool. It is intended for educational purposes only. Users are responsible for ensuring lawful usage.

## Table of Contents
- [SugarFree](#sugarfree)
  - [Description](#description)
  - [Disclaimer](#disclaimer)
  - [Table of Contents](#table-of-contents)
  - [Acknowledgement](#acknowledgement)
  - [Installation](#installation)
  - [Usage](#usage)
  - [References](#references)

## Acknowledgement

This project created with :heart: by [@nickvourd](https://x.com/nickvourd) && [@IAMCOMPROMISED](https://x.com/IAMCOMPROMISED).

Special thanks to my friend [Marios Gyftos](https://www.linkedin.com/in/marios-gyftos-a6b62122/) for inspiring the concept of automated stages.

## Installation

You can use the [precompiled binaries](https://github.com/nickvourd/SugarFree/releases), or you can manually install SugarFree by following the next steps:

1) Clone the repository by executing the following command:

```
git clone https://github.com/nickvourd/SugarFree.git
```

2) Once the repository is cloned, navigate into the SugarFree directory:

```
cd SugarFree
```

3) Install the third-party dependencies:

```
go mod download
```

4) Build SugarFree with the following command:

```
go build SugarFree
```

## Usage

:information_source: Please refer to the [SugarFree Wiki](https://github.com/nickvourd/SugarFree/wiki) for detailed usage instructions and examples of commands.

```
███████╗██╗   ██╗ ██████╗  █████╗ ██████╗ ███████╗██████╗ ███████╗███████╗
██╔════╝██║   ██║██╔════╝ ██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔════╝██╔════╝
███████╗██║   ██║██║  ███╗███████║██████╔╝█████╗  ██████╔╝█████╗  █████╗  
╚════██║██║   ██║██║   ██║██╔══██║██╔══██╗██╔══╝  ██╔══██╗██╔══╝  ██╔══╝  
███████║╚██████╔╝╚██████╔╝██║  ██║██║  ██║██║     ██║  ██║███████╗███████╗
╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚══════╝╚══════╝                                                        

SugarFree v2.0 - Less sugar (entropy) for your binaries.
SugarFree is an open source tool licensed under MIT.
Written with <3 by @nickvourd && @IAMCOMPROMISED...
Please visit https://github.com/nickvourd/SugarFree for more...

Usage:
  SugarFree [flags]
  SugarFree [command]

Aliases:
  SugarFree, sugarfree, SUGARFREE, sf

Available Commands:
  free        Free command
  help        Help about any command
  info        Info command

Flags:
  -h, --help      help for SugarFree
  -v, --version   Show SugarFree current version
```

## References

- [Threat Hunting with File Entropy by Practical Security Analytics LLC](https://practicalsecurityanalytics.com/file-entropy/)
- [Using Entropy in Threat Hunting: a Mathematical Search for the Unknown by Red Canary](https://redcanary.com/blog/threat-detection/threat-hunting-entropy/)
- [EntropyReducer GitHub by Maldev Academy](https://github.com/Maldev-Academy/EntropyReducer)
- [EntropyFix GitHub by EspressoCake](https://github.com/EspressoCake/EntropyFix)
