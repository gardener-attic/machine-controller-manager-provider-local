# machine-controller-manager-provider-local
[![REUSE status](https://api.reuse.software/badge/github.com/gardener/machine-controller-manager-provider-local)](https://api.reuse.software/info/github.com/gardener/machine-controller-manager-provider-local)
[![Go Report Card](https://goreportcard.com/badge/github.com/gardener/machine-controller-manager-provider-local)](https://goreportcard.com/report/github.com/gardener/machine-controller-manager-provider-local)

Out of tree (controller-based) implementation for `local` as a new provider.
The local out-of-tree provider implements the interface defined at [MCM OOT driver](https://github.com/gardener/machine-controller-manager/blob/master/pkg/util/provider/driver/driver.go).

## Fundamental Design Principles

Following are the basic principles kept in mind while developing the external plugin.

- Communication between this Machine Controller (MC) and Machine Controller Manager (MCM) is achieved using the Kubernetes native declarative approach.
- Machine Controller (MC) behaves as the controller used to interact with the cloud provider AWS and manage the VMs corresponding to the machine objects.
- Machine Controller Manager (MCM) deals with higher level objects such as machine-set and machine-deployment objects.

## Testing the Controller

1. Open terminal and change directory to `$GOPATH/src/github.com/gardener`. Clone this repository.

2. Navigate to `$GOPATH/src/github.com/gardener/machine-controller-manager-provider-local`:
    - In the `MAKEFILE` make sure `$TARGET_KUBECONFIG` points to the cluster where you wish to manage machines. `$CONTROL_NAMESPACE` represents the namespaces where MCM is looking for machine CR objects, and `$CONTROL_KUBECONFIG` points to the cluster which holds these machine CRs.
    - Run the machine controller (driver) using the command below.

        ```bash
        make start
        ```

3. On the second terminal pointing to `$GOPATH/src/github.com/gardener`,
    - Clone the [latest MCM code](https://github.com/gardener/machine-controller-manager):

        ```bash
        git clone git@github.com:gardener/machine-controller-manager.git
        ```

    - Navigate to the newly created directory:

        ```bash
        cd machine-controller-manager
        ```

    - Deploy the required CRDs from the machine-controller-manager repo:

        ```bash
        kubectl apply -f kubernetes/crds.yaml
        ```

    - Run the machine-controller-manager:

        ```bash
        make start
        ```

4. On the third terminal pointing to `$GOPATH/src/github.com/gardener/machine-controller-manager-provider-local`
    - Fill in the object files given below and deploy them as described below.
    - Deploy the `machine-class`

        ```bash
        kubectl apply -f kubernetes/machine-class.yaml
        ```

    - Deploy the `kubernetes secret` if required.

        ```bash
        kubectl apply -f kubernetes/secret.yaml
        ```

    - Deploy the `machine` object and make sure it joins the cluster successfully.

        ```bash
        kubectl apply -f kubernetes/machine.yaml
        ```

    - Once machine joins, you can test by deploying a machine-deployment.

    - Deploy the `machine-deployment` object and make sure it joins the cluster successfully.

        ```bash
        kubectl apply -f kubernetes/machine-deployment.yaml
        ```

    - Make sure to delete both the `machine` and `machine-deployment` object after use.

        ```bash
        kubectl delete -f kubernetes/machine.yaml
        kubectl delete -f kubernetes/machine-deployment.yaml
        ```

Static code checks and tests can be executed by running `make verify`. We are using Go modules for Golang package dependency management and [Ginkgo](https://github.com/onsi/ginkgo)/[Gomega](https://github.com/onsi/gomega) for testing.

## Feedback and Support

Feedback and contributions are always welcome. Please report bugs or suggestions as [GitHub issues](https://github.com/gardener/machine-controller-manager-provider-local/issues) or join our [Slack channel #gardener](https://kubernetes.slack.com/messages/gardener) (please invite yourself to the Kubernetes workspace [here](http://slack.k8s.io)).
