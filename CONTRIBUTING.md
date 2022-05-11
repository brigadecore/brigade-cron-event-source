# Contributing Guide

The Brigade Cron Event Source is an official extension of the Brigade project
and as such follows all of the practices and policies laid out in the main
[Brigade Contributor Guide](https://docs.brigade.sh/topics/contributor-guide/).
Anyone interested in contributing to this gateway should familiarize themselves
with that guide _first_.

The remainder of _this_ document only supplements the above with things specific
to this project.

## Running `make hack-kind-up`

As with the main Brigade repository, running `make hack-kind-up` in this
repository will utilize [ctlptl](https://github.com/tilt-dev/ctlptl) and
[KinD](https://kind.sigs.k8s.io/) to launch a local, development-grade
Kubernetes cluster that is also connected to a local Docker registry.

In contrast to the main Brigade repo, this cluster is not pre-configured for
building and running Brigade itself from source, rather it is pre-configured for
building and running _this event source_ from source. Because Brigade is a
logical prerequisite for this event source to be useful in any way, `make
hack-kind-up` will pre-install a recent, _stable_ release of Brigade into the
cluster.

## Running `tilt up`

As with the main Brigade repository, running `tilt up` will build and deploy
project code (the event source, in this case) from source.

For the event source to successfully communicate with the Brigade instance in
your local, development-grade cluster, you will need to execute the following
steps _before_ running `tilt up`:

1. Log into Brigade:

   ```shell
   $ brig login -k -s https://localhost:31600 --root
   ```

   The root password is `F00Bar!!!`.

1. Create a service account for the gateway:

   ```shell
   $ brig service-account create \
       --id cron-event-source \
       --description cron-event-source
   ```

1. Copy the token returned from the previous step and export it as the
   `BRIGADE_API_TOKEN` environment variable:

   ```shell
   $ export BRIGADE_API_TOKEN=<token from previous step>
   ```

1. Grant the service account permission to create events:

   ```shell
   $ brig role grant EVENT_CREATOR \
     --service-account cron-event-source \
     --source brigade.sh/cron
   ```

   > ⚠️&nbsp;&nbsp;Contributions that automate the creation and configuration of
   > the service account setup are welcome.

1. Edit the `events` section of `charts/brigade-cron-event-source/values.yaml`
   to describe the events you'd like to emit into Brigade and the schedule on
   which they should be emitted. Refer to instructions in the
   [README](README.md) for more information.

   > ⚠️&nbsp;&nbsp;Take care not to include modifications to the `values.yaml`
   > file in any PRs you open.

You can then run `tilt up` to build and deploy this gateway from source.
