module github.com/gardener/machine-controller-manager-provider-local

go 1.16

require (
	github.com/gardener/gardener v1.36.0
	github.com/gardener/machine-controller-manager v0.41.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.1.7
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/component-base v0.22.2
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/controller-runtime v0.10.2
)

replace (
	k8s.io/api => k8s.io/api v0.20.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.6
	k8s.io/client-go => k8s.io/client-go v0.20.6
	k8s.io/component-base => k8s.io/component-base v0.20.6
	k8s.io/klog => k8s.io/klog v0.4.0
	k8s.io/utils => k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.1
)
