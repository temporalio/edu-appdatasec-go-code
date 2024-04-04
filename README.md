# Code Repository for Securing Application Data in Temporal (Go)
This repository provides code used for exercises and demonstrations
included in the Go version of the 
[Securing Application Data in Temporal](https://learn.temporal.io/courses/appdatasec) 
training course.

It's important to remember that the example code used in this course was designed to support learning a specific aspect of Temporal, not to serve as a ready-to-use template for implementing a production system.

## Hands-On Exercises

Directory Name                     | Exercise
:--------------------------------- | :-------------------------------------------------------
`exercises/custom-converter`       | [Exercise 1](exercises/custom-converter/README.md)
`exercises/codec-server`           | [Exercise 2](exercises/codec-server/README.md)

Samples
| Directory Name                     | Description                                                                                                                                     |
| :--------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------- |
| `samples/composite-converter` | [Sample of how to provide an additional Payload Converter to your Composite Converter](samples/composite-converter/README.md) |

## Reference
The following links provide additional information that you may find helpful as you work through this course.
* [General Temporal Documentation](https://docs.temporal.io/)
* [Temporal Go SDK Documentation](https://pkg.go.dev/go.temporal.io/sdk)
* [Go Language Documentation](https://go.dev/doc/)
* [GitPod Documentation: Troubleshooting](https://www.gitpod.io/docs/troubleshooting)


## Exercise Environment for this Course
You can launch an exercise environment for this course in GitPod by 
clicking the button below:

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/temporalio/edu-appdatasec-go-code)

Alternatively, you can follow 
[these instructions](https://learn.temporal.io/getting_started/go/dev_environment/) to 
set up your own Temporal Cluster with Docker Compose, which you can use as an 
exercise environment.
