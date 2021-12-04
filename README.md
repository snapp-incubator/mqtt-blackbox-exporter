# MQTT Blackbox Exporter

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/snapp-incubator/mqtt-blackbox-exporter/ci?label=ci&logo=github&style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/snapp-incubator/mqtt-blackbox-exporter.svg)](https://pkg.go.dev/github.com/snapp-incubator/mqtt-blackbox-exporter)
[![Codecov](https://img.shields.io/codecov/c/gh/snapp-incubator/mqtt-blackbox-exporter?logo=codecov&style=flat-square)](https://codecov.io/gh/snapp-incubator/mqtt-blackbox-exporter)

## Introduction

In each probe it sends a message over MQTT broker and then wait for getting it over subscription.
By measuring this time and also connection, subscription and etc durations you can check your cluster status.
