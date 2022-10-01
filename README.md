# ruleengine-core

[![Build](https://github.com/niharrathod/ruleengine-core/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/niharrathod/ruleengine-core/actions/workflows/build.yml)
[![Test](https://github.com/niharrathod/ruleengine-core/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/niharrathod/ruleengine-core/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/niharrathod/ruleengine-core/badge.svg?branch=master)](https://coveralls.io/github/niharrathod/ruleengine-core?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Overview

ruleengine-core is a strictly typed rule engine library, provides a simple interface to create ruleengine and evaluate rules for given input. It decouples and abstract-out logic for rule setup and rule evaluation. You can store rules(RuleEngineConfig) as a JSON in a store outside of core rule engine logic, which makes it easy for rule engine versioning.

## Example
