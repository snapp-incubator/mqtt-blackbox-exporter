<h1 align="center"> MQTT Blackbox Exporter </h1>

<p align="center">
    <img src="./.github/assets/logo.jpg" height="250px">
</p>

<p align="center">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/snapp-incubator/mqtt-blackbox-exporter/ci.yaml?logo=github&style=for-the-badge">
    <img alt="Codecov" src="https://img.shields.io/codecov/c/github/snapp-incubator/mqtt-blackbox-exporter?logo=codecov&style=for-the-badge">
    <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/snapp-incubator/mqtt-blackbox-exporter?logo=github&style=for-the-badge">
    <img alt="GitHub tag (with filter)" src="https://img.shields.io/github/v/tag/snapp-incubator/mqtt-blackbox-exporter?style=for-the-badge&logo=git">
    <img alt="GitHub go.mod Go version (subdirectory of monorepo)" src="https://img.shields.io/github/go-mod/go-version/snapp-incubator/mqtt-blackbox-exporter?style=for-the-badge&logo=go">
</p>

## Introduction

In each probe it sends a message over MQTT broker and then wait for getting it over subscription.
By measuring this time and also connection, subscription etc. durations you can check your cluster status.
At Snapp! we use it to detect our EMQX clusters status from the client perspective.
