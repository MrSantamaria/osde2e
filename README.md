# OSDe2e

[![GoDoc](https://godoc.org/github.com/openshift/osde2e?status.svg)](https://godoc.org/github.com/openshift/osde2e)

## Introduction

A comprehensive test framework used for Service Delivery to test all aspects of
Managed OpenShift Clusters ([OpenShift Dedicated]). The data generated by
the different test coverage is used to inform product releases and decisions.

OSDe2e key features are:

* Portable test framework that can run anywhere to validate end to end test workflows
  * Run locally from a developers workstation or from a CI application
* Supports create/delete different cluster deployment types
  * ROSA, ROSA Hosted Control Plane (e.g. HyperShift), OSD
* Performs cluster health checks to ensure cluster is operational prior to
  running tests
* Perform cluster upgrades
* Captures artifacts such as logs, metrics, metadata to be archived for later usage
* Tests OSD operators along with other OpenShift features from a
  customer/SRE point of view
* Provides a test harness to validate [Add Ons][OSDE2E Test Harness]

When osde2e is invoked, the standard workflow is followed:

* Load configuration
* Cluster deployment (when not leveraging an existing cluster)
* Verify the health of the cluster
* Run tests (pre upgrade)
* Collect logs, metrics and metadata
* Upgrade cluster (when defined)
* Verify the health of the cluster post upgrade
* Run tests (post upgrade - when upgrade is defined)
* Collect logs, metrics and metadata
* Cluster deprovision (when this is toggled on)

## Prerequisites

Prior to running osde2e, make sure you meet the minimal prerequisites defined below:

* Navigate to [OpenShift Cluster Manager (OCM)][OpenShift Offline Token] to obtain
  an OpenShift offline token.
  * Save your token into the environment variable `OCM_TOKEN` for later usage
* Verify (submit a request if required) your Red Hat account has adequate quota for
  deploying clusters based on your preferred deployment type
* A go workspace running the minimal version defined in the [go.mod](go.mod)

## Run

OSDe2e can be invoked by one of two ways. Refer to each section below to learn how
to run it.

### From Source

Running from source requires you to build the osde2e binary. Follow the steps below
to do this:

```shell
git clone https://github.com/openshift/osde2e.git
cd osde2e
go mod tidy
make build
```

On completion of the `make build` target, the generated binary will reside in the
directory `./out/`. Where you can then invoke osde2e `./out/osde2e --help`.

### From Container Image

Running from a container image using a container engine (e.g. docker, podman).
You can either build the image locally or consume the public image available on
[quay.io][OSDE2E Quay Image].

```shell
export CONTAINER_ENGINE=<docker|podman>

# Build Image
make build-image
$CONTAINER_ENGINE run quay.io/app-sre/osde2e:latest <args>

# Pull Image
$CONTAINER_ENGINE pull quay.io/app-sre/osde2e:latest
$CONTAINER_ENGINE run quay.io/app-sre/osde2e:latest <args>
```

## Config Input

OSDe2e provides multiple ways for you to provide input to tailor what test workflows
you wish to validate. It provides four ways for you to provide input
(order is lowest to highest precedence):

* Use pre-canned composable default [configs]
* Provide a custom config
* Environment variables
* Command line options

*It is highly recommended to leave sensitive settings as environment variables (e.g. `OCM_TOKEN`).
This way the chance of these settings defined in a custom config file are not checked into
source control.*

### Pre-Canned Default Configs

The [configs] package provides pre-canned default configs available for you to use.
These config files are named based on what action they are performing. Within the config
file can contain multiple settings to tailor osde2e.

Example config [stage](configs/stage.yaml):

This default config is telling osde2e to use the stage OCM environment.

```yaml
ocm:
  env: stage
```

You can provide N+1 pre-canned configs to osde2e. Example below will deploy
a OSD cluster within the OCM stage environment.

```shell
./out/osde2e test --configs aws,stage
```

### Custom Config

The composable configs consist of a number of small YAML files that can all be loaded together.
Rather than using the built in configs, you can also elect to build your own custom YAML file
and provide that using the `--custom-config` CLI option.

```shell
osde2e test --custom-config ./osde2e.yaml
```

The custom config below is a basic example for deploying a ROSA STS cluster and running
all of the OSD operators tests that do not have the informing label associated to them.

```yaml
dryRun: false
provider: rosa
cloudProvider:
  providerId: aws
  region: us-east-1
rosa:
  env: stage
  STS: true
cluster:
  name: osde2e
tests:
  ginkgoLabelFilter: Operators && !Informing
```

You can use both pre-canned default configs and your own custom configs:

```shell
./out/osde2e test --configs aws --custom-config ./osde2e.yaml
```

### Environment Variables

Any config option can be passed in using environment variables.
Please refer to the [config package] for exact environment variable names.

Below is an example to spin up a OSD cluster and test it:

```shell
OCM_TOKEN=<ocm-token> \
OSD_ENV=prod \
CLUSTER_NAME=my-cluster \
MAJOR_TARGET=4 \
MINOR_TARGET=12 \
osde2e test
```

These also can be combined with pre-canned default configs and custom configs:

```shell
OCM_TOKEN=<ocm-token> \
MAJOR_TARGET=4 \
MINOR_TARGET=12 \
osde2e test --configs prod,e2e-suite
```

```shell
OCM_TOKEN=<ocm-token> \
MAJOR_TARGET=4 \
MINOR_TARGET=12 \
osde2e test --configs prod,e2e-suite
```

A list of commonly used environment variables are included in [Config variables].

### Command Line Options

Some configuration settings are also exposed as command-line options.
A full list can be displayed by providing `--help` after the command.

Below is an example of using options for the `test` command:

```shell
./out/osde2e test --cluster-id <cluster-id> \
  --provider stage \
  --skip-health-check \
  --focus-tests "RBAC Operator"
```

Another example below is you can skip cluster health check, must gather
as follows.

```shell
POLLING_TIMEOUT=1 \
./out/osde2e test --cluster-id=<cluster-id> \
--configs stage \
--must-gather=False \
--skip-health-check \
--focus-tests="rh-api-lb-test"
```

A list of commonly used CLI flags are included in [Config variables].

### Examples

To see more examples of configuring input for osde2e, refer to the
[prowgen jobs][OSDE2E ProwGen Job Config] in the OpenShift release repository
owned by the team. These will be always up to date with the latest changes
osde2e has to offer.

## Cluster Deployments

OSDe2e provides native support for deploying the following cluster types:

* ROSA
* ROSA Hosted Control Plane (HyperShift)
* OSD (OpenShift Dedicated)

You can have osde2e deploy the cluster if a cluster ID is not provided or
you can leverage an existing cluster by giving the cluster ID as input at
runtime.

You can also provide it a kubeconfig file and osde2e can attempt to target
that cluster.

```shell
export TEST_KUBECONFIG=<kubeconfig-file>
./out/osde2e test <args>
```

*It may be possible to test against a non Managed OpenShift cluster
(a traditional OpenShift Container Platform cluster). Though this will
require you to alter the input settings as non managed clusters will not
have certain items applied to them like a Managed cluster would (e.g. OSD
operators, health checks, etc).*

## Tests

OSDe2e currently holds all core and operator specific tests and are maintained by the CICD team.
Test types range from core OSD verification, OSD operators to scale/conformance.

*Currently in flight: OSD operator tests will no longer reside in osde2e repository and
live directly alongside the operator source code in its repository*

### Selecting Tests To Run

OSDe2e supports a couple different ways you can select which tests you would like to run. Below presents
the commonly used methods for this:

* Using the label filter. Labels are ginkgos way to tag test cases. The examples below
   will tell osde2e to run all tests that have the `E2E` label applied.

```shell
# Command line option
osde2e test --label-filter E2E

# Passed in using a custom config file
tests:
  ginkgoLabelFilter: E2E
```

* Using focus strings. Focus strings are ginkgos way to select test cases based on string regex.

```shell
# Command line option
osde2e test --focus-tests "OCM Agent Operator"

# Custom config file
tests:
  focus: "OCM Agent Operator"
```

* Using a combination of labels and focus strings to fine tune your test selection.
   The examples below tell osde2e to run all ocm agent operator tests and avoid running
   the upgrade test case.

```shell
# Command line options
osde2e test --label-filter "Operators && !Upgrade" --focus-tests "OCM Agent Operator"

# Custom config file
tests:
  ginkgoLabelFilter: "Operators && !Upgrade"
  focus: "OCM Agent Operator"
```

### Writing Tests

Refer to the [Writing Tests] document for guidelines and standards.

Third-party (Addon) tests are built as containers that spin up and report back results to OSDe2e.
These containers are built and maintained by external groups looking to get CI signal for
their product within OSD. The definition of a third-party test is maintained within
the `managed-tenants` repo and is returned via the Add-Ons API.

For more information please see the [OSDE2E Test Harness] repository to learn more
for writing add on tests.

## Reporting

Each time osde2e runs it captures as much data that it possible can. Data can include
cluster/pod logs, prometheus metrics, test data generated, hive version and osde2e version
to identify any possible flakiness in the environment.

Each time tests are executed a JUnit XML file will be generated to capture all the tests
that ran and statistics about them (e.g. pass/fail, duration). These XML files will be later
used by external applications to present metrics and data for others to see into. An example of
this is they are used to present data in [TestGrid Dashboards][TestGrid Dashboard].

### CI/CD Job Results Database

We have provisioned an AWS RDS Postgres database to store information about our CI jobs
and the tests that they execute. We used to store our data only within prometheus,
but prometheus's timeseries paradigm prevented us from being able to express certain
queries (even simple ones like "when was the last time this test failed").

The test results database (at time of writing) stores data about each job and its configuration,
as well as about each test case reported by the Junit XML output of the job.

This data allows us to answer questions about frequency of job/test failure,
relationships between failures, and more. The code responsible for managing the
database can be found in the [`./pkg/db/`](pkg/db) directory,
along with a README describing how to develop against it.

### Database usage from OSDe2e

Because `osde2e` runs a a cluster of parallel, ephemeral prow jobs, our usage of
the database is unconventional. We have to write all of our database interaction
logic with the understanding that any number of other prow jobs could be modifying
the data at the same time that we are.

We use the database to generate alerts for the CI Watcher to use, and we follow
this algorithm to generate those alerts safely in our highly-concurrent usecase
(at time of writing, implemented [here](https://github.com/openshift/osde2e/blob/cfd38c75532274d619840ad505c1232881eb417a/pkg/e2e/e2e.go#L1029)):

1. At the end of each job, list all testcases that failed during the current job. Implemented by [`ListAlertableFailuresForJob`](https://github.com/openshift/osde2e/blob/cfd38c75532274d619840ad505c1232881eb417a/pkg/db/queries/queries.sql#L66).
1. Generate a list of testcases (in any job) that have failed more than once during the last 48 hours. Implemented by [`ListProblematicTests`](https://github.com/openshift/osde2e/blob/cfd38c75532274d619840ad505c1232881eb417a/pkg/db/queries/queries.sql#L105).
1. For each failing testcase in the current job, create a PD alert if the testcase is one of those that have failed more than once in the last 48 hours.
1. After generating all alerts as above, merge all pagerduty alerts that indicate failures for the same testcase (this merge uses the title of the alert, which is the testcase name, to group the alerts).
1. Finally, close any PD incident for a testcase that does not appear in the list of testcases failing during the last 48 hours.

> Why does each job only report its own failures? The database is global, and a single job could report for all of them.

If each job reported for the failures of all recent jobs, we'd create an enormous number of redundant alerts for no benefit. Having each job only report its own failures keeps the level of noise down *without* requiring us to build some kind of concensus mechanism between the jobs executing in parallel.

> Why close the PD incidents for test cases that haven't failed in the last 48 hours?

This is a heuristic designed to automatically close incidents when the underlying test problem has been dealt with. If we stop seeing failures for a testcase, it probably means that the testcase has stopped failing. This can backfire, and a more intelligent heuristic is certainly possible.

[Config variables]:/docs/Config.md
[configs]:/configs/
[config package]:/pkg/common/config/config.go
[OSDE2E Quay Image]: quay.io/app-sre/osde2e
[OpenShift Dedicated]: https://docs.openshift.com/dedicated/welcome/index.html
[OSDE2E Test Harness]: https://github.com/openshift/osde2e-example-test-harness
[OpenShift Offline Token]:https://cloud.redhat.com/openshift/token
[OSDE2E ProwGen Job Config]: https://github.com/openshift/release/blob/master/ci-operator/config/openshift/osde2e/openshift-osde2e-main.yaml
[TestGrid Dashboard]: https://testgrid.k8s.io/redhat-openshift-osd
[Writing Tests]:/docs/Writing-Tests.md
