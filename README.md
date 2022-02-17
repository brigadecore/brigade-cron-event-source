# Brigade Cron Event Source

![build](https://badgr.brigade2.io/v1/github/checks/brigadecore/brigade-cron-event-source/badge.svg?appID=99005)
[![codecov](https://codecov.io/gh/brigadecore/brigade-cron-event-source/branch/main/graph/badge.svg?token=ZSac3nWz6M)](https://codecov.io/gh/brigadecore/brigade-cron-event-source)
[![Go Report Card](https://goreportcard.com/badge/github.com/brigadecore/brigade-cron-event-source)](https://goreportcard.com/report/github.com/brigadecore/brigade-cron-event-source)
[![slack](https://img.shields.io/badge/slack-brigade-brightgreen.svg?logo=slack)](https://kubernetes.slack.com/messages/C87MF1RFD)

<img width="100" align="left" src="logo.png">

The Brigade Cron Event Source offers an easy, low-overhead method of emitting
user-defined events into Brigade's event bus on a user-defined schedule. It will
create a Kubernetes `CronJob` resource for each such event. A small program will
wake on the specified schedule to emit its event. _That's it._

<br clear="left"/>

## Why a Cron Event Source?

Before getting started with this event source, do ask yourself if you truly need
it. If your only interest is in executing a _simple_ task on a particular
schedule, Kubernetes, by itself, gives you everything you need. In such a case,
we don't recommend over-complicating things by involving Brigade!

If, however, you wish to execute a more complex workflow on a particular
schedule or if your use case is well-served by other Brigade features, _do_
consider utilizing this event source.

A non-exhaustive set of reasons you might wish to use this gateway includes:

* Your workflow is complex and involves multiple containers that need to execute
  concurrently or in serial.
* You also have _other_ events you need to handle.
* You care about capturing and saving logs from your automated task.

## Installation

Prerequisites:

* A Kubernetes cluster:
    * For which you have the `admin` cluster role
    * That is already running Brigade 2

* `kubectl`, `helm` (commands below require Helm 3.7.0+), and `brig` (the
  Brigade 2 CLI)

### 1. Create a Service Account for the Gateway

> ⚠️&nbsp;&nbsp;To proceed beyond this point, you'll need to be logged into Brigade 2
as the "root" user (not recommended) or (preferably) as a user with the `ADMIN`
role. Further discussion of this is beyond the scope of this documentation.
Please refer to Brigade's own documentation.

Using Brigade 2's `brig` CLI, create a service account for the event source to
use:

```console
$ brig service-account create \
    --id brigade-cron-event-source \
    --description brigade-cron-event-source
```

Make note of the __token__ returned. This value will be used in another step.
_It is your only opportunity to access this value, as Brigade does not save it._

Authorize this service account to create new events:

```console
$ brig role grant EVENT_CREATOR \
    --service-account brigade-cron-event-source \
    --source brigade.sh/cron
```

> ⚠️&nbsp;&nbsp;The `--source brigade.sh/cron` option specifies that this service
account can be used _only_ to create events having a value of `brigade.sh/cron`
in the event's `source` field. _This is a security measure that prevents the
event source from using this token for impersonating other sources._

### 2. Install the Cron Event Source

> ⚠️&nbsp;&nbsp;be sure you are using
> [Helm 3.7.0](https://github.com/helm/helm/releases/tag/v3.7.0) or greater and
> enable experimental OCI support:
>
> ```console
>  $ export HELM_EXPERIMENTAL_OCI=1
>  ```

As this chart requires custom configuration to function properly, we'll need to
create a chart values file with said config.

Use the following command to extract the full set of configuration options into
a file you can modify:

```console
$ helm inspect values oci://ghcr.io/brigadecore/brigade-cron-event-source \
    --version v0.1.1 > ~/brigade-cron-event-source-values.yaml
```

Edit `~/brigade-cron-event-source-values.yaml`, making the following changes:

* `brigade.apiAddress`: Address of the Brigade API server, beginning with
  `https://`

* `brigade.apiToken`: Service account token from step 2

* `events`: By default, this field contains an empty array. It can be modified
  to enumerate events that should each be emitted to Brigade's event bus on a
  schedule of your choosing. The file contains extensive comments on how to do
  this, but some highlights are covered here:

  * `name`: Give your event a name. It must be unique to this installation of
    the cron event source. The name may contain only lower case alphanumeric
    characters and dashes.

  * `schedule`: The schedule for emitting the event to Brigade. It should follow
    the syntax described 
    [here](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#schedule).

  * `brigadeEvent`: YAML representation of a Brigade event. It must conform to
    the [v2 event schema](https://schemas.brigade.sh/schemas-v2/event.json) and
    has the following _additional_ restrictions:

      * `projectID`: MUST be specified

      * `source`: MUST be `brigade.sh/cron`

Save your changes to `~/brigade-cron-event-source-values.yaml` and use the
following command to install the event source using the above customizations:

```console
$ helm install brigade-cron-event-source \
    oci://ghcr.io/brigadecore/brigade-cron-event-source \
    --version v0.1.1 \
    --create-namespace \
    --namespace brigade-cron-event-source \
    --values ~/brigade-cron-event-source-values.yaml \
    --wait \
    --timeout 300s
```

### 3. Subscribe a Brigade Project

As noted in the previous section, events emitted by this event source MUST
specify a Brigade project by ID. For each event emitted by this event source,
create or modify the corresponding project.

In the example (new) project definition below, we subscribe to an event of type
`cleanup-requested` emitted by this event source:

```yaml
apiVersion: brigade.sh/v2
kind: Project
metadata:
  id: cron-demo
description: A project that demonstrates integration with the cron event source
spec:
  eventSubscriptions:
  - source: brigade.sh/cron
    types:
    - cleanup-requested
  workerTemplate:
    defaultConfigFiles:
      brigade.js: |-
        const { events } = require("@brigadecore/brigadier");

        events.on("brigade.sh/cron", "cleanup-requested", () => {
          console.log("Performing nightly cleanup...");
          // ...
        });

        events.process();
```

Assuming this file were named `project.yaml`, you can create the project like
so:

```console
$ brig project create --file project.yaml
```

## General Guidance

This event source is different than most others in that it has the flexibility
to generate user-defined events on a user-defined schedule. With this being the
case, this event source's creators wish to offer some guidance on how best to
leverage that flexibility without becoming encumbered by it.

In general, what we do _not_ recommend is emitting events with `type` values
that reflect times or intervals. `midnight` or `top-of-the-hour`, are in our
opinion, poor event types. While these might accurately convey real world
events, such as "the stroke of midnight," such events can be problematic for
their corresponding projects. Supposing, for instance, that some project
executes a nightly cleanup process in response to the `midnight` event and one
wished to move nightly cleanup to 1:00 AM. Not only would the event source need
to be reconfiguration to emit events with a `type` value like (for instance)
`one-am`, but the project's subscriptions and script would _likewise_ require
updating.

In contrast, what we _do_ recommend is emitting events with `type` values that
reflect some _desired effect_. If the event source were configured to emit
events with `type` value `cleanup-requested`, for instance, then rescheduling
the event is only a matter of updating event source configuration and neither
the corresponding project's event subscriptions nor script need to be modified.

## Contributing

The Brigade project accepts contributions via GitHub pull requests. The
[Contributing](CONTRIBUTING.md) document outlines the process to help get your
contribution accepted.

## Support & Feedback

We have a Slack channel! Visit [slack.brigade.sh](https://slack.brigade.sh) to
join us. We welcome any support questions or feedback.

To report an issue or to request a feature, open an issue
[here](https://github.com/brigadecore/brigade-cron-event-source/issues).

## Code of Conduct

Participation in the Brigade project is governed by the
[CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).
