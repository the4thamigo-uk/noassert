# noassert
_Static analysis tool to detect assertion-free tests in golang projects_

[![Build Status](https://travis-ci.org/the4thamigo-uk/noassert.svg?branch=master)](https://travis-ci.org/the4thamigo-uk/noassert)
[![Coverage Status](https://coveralls.io/repos/github/the4thamigo-uk/noassert/badge.svg?branch=master)](https://coveralls.io/github/the4thamigo-uk/noassert?branch=master)

## Description

To write useful unit tests, developers must provide assertions to validate behaviour is correct. However, developers can sometimes forget to add such assertions, or worse, 
they might write [assertion-free tests](https://martinfowler.com/bliki/AssertionFreeTesting.html) in order to obtain good code coverage. 

This simple static-analysis tool aims to flag up any tests that do not call `testing.Fail` at some point in their callgraph.

The tool does nothing to help you write 'good' unit tests however...

## Getting Started

Run :

```
go get github.com/the4thamigo_uk/noassert
```

This will build and install the binary in the GOROOT/bin folder.

To use simply specify a package path e.g. :

```
noassert "github.com/the4thamigo-uk/noassert/testlib"
```

## Road Map

Improvements will be made to support additional command line arguments.
