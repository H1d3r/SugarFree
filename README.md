# SugarFree

Less sugar (entropy) for your binaries

<p align="center">
  <img width="500" height="400" src="/Pictures/logo3.png"><br /><br />
  <!--<img alt="GitHub License" src="https://img.shields.io/github/license/nickvourd/SugarFree?style=social&logo=GitHub&logoColor=purple">
  <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/nickvourd/SugarFree?logoColor=yellow">
  <img alt="GitHub forks" src="https://img.shields.io/github/forks/nickvourd/SugarFree?logoColor=red">
  <img alt="GitHub watchers" src="https://img.shields.io/github/watchers/nickvourd/SugarFree?logoColor=blue">
  <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/nickvourd/SugarFree?style=social&logo=GitHub&logoColor=green">-->
</p>

## Description

SugarFree is an open-source tool designed to analyze and reduce the entropy of a provided PE file. SugarFree appends null bytes (`0x00`) to the end of the file, increasing the binary size while reducing its entropy.

The following list explains the meaning of each SugarFree command:

- **info**: Calculates the entropy of a PE file and its sections.
- **free**: Lowers the overall entropy of a PE file.

SugarFree is written in Golang, a cross-platform language, enabling its use on both Windows and Linux systems.

This project created with :heart: by [@nickvourd](https://x.com/nickvourd) && [@IAMCOMPROMISED](https://x.com/IAMCOMPROMISED).

Special thanks to my friend [Marios Gyftos](https://www.linkedin.com/in/marios-gyftos-a6b62122/) for inspiring the concept of automated stages.

## Table of Contents
- [SugarFree](#sugarfree)
  - [Description](#description)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
  - [References](#references)

## Installation

