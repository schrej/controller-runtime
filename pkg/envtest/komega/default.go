package komega

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Obj[V any] interface {
	*V
	client.Object
}

type ObjList[V any] interface {
	*V
	client.ObjectList
}

type GenericUpdateFunc[T client.Object] func(T client.Object)

// defaultK is the Komega used by the package global functions.
var defaultK = &komega{}

// SetDefaultClient sets the client used by the package global functions.
func SetClient(c client.Client) {
	defaultK = &komega{client: c}
}

func checkClient() {
	if defaultK.client == nil {
		panic("Komega's client is not set. Use SetClient to set it.")
	}
}

// Get returns a function that fetches a resource and returns the occurring error.
// It can be used with gomega.Eventually() like this
//   deployment := appsv1.Deployment{ ... }
//   gomega.Eventually(komega.Get(&deployment)).To(gomega.Succeed())
// By calling the returned function directly it can also be used with gomega.Expect(komega.Get(...)()).To(...)
func Get(obj client.Object) func() error {
	checkClient()
	return defaultK.Get(obj)
}

// List returns a function that lists resources and returns the occurring error.
// It can be used with gomega.Eventually() like this
//   deployments := v1.DeploymentList{ ... }
//   gomega.Eventually(k.List(&deployments)).To(gomega.Succeed())
// By calling the returned function directly it can also be used as gomega.Expect(k.List(...)()).To(...)
func List(obj client.ObjectList, opts ...client.ListOption) func() error {
	checkClient()
	return defaultK.List(obj, opts...)
}

// Update returns a function that fetches a resource, applies the provided update function and then updates the resource.
// It can be used with gomega.Eventually() like this:
//   deployment := appsv1.Deployment{ ... }
//   gomega.Eventually(k.Update(&deployment, func (o client.Object) {
//     deployment.Spec.Replicas = 3
//     return &deployment
//   })).To(gomega.Scucceed())
// By calling the returned function directly it can also be used as gomega.Expect(k.Update(...)()).To(...)
func Update[T Obj[V], V any](nn types.NamespacedName, f GenericUpdateFunc[T], opts ...client.UpdateOption) func() error {
	checkClient()
	var obj T = new(V)
	obj.SetName(nn.Name)
	obj.SetNamespace(nn.Namespace)
	return defaultK.Update(obj, func() { f(obj) }, opts...)
}

// UpdateStatus returns a function that fetches a resource, applies the provided update function and then updates the resource's status.
// It can be used with gomega.Eventually() like this:
//   deployment := appsv1.Deployment{ ... }
//   gomega.Eventually(k.Update(&deployment, func (o client.Object) {
//     deployment.Status.AvailableReplicas = 1
//     return &deployment
//   })).To(gomega.Scucceed())
// By calling the returned function directly it can also be used as gomega.Expect(k.UpdateStatus(...)()).To(...)
func UpdateStatus[T Obj[V], V any](nn types.NamespacedName, f GenericUpdateFunc[T], opts ...client.UpdateOption) func() error {
	checkClient()
	var obj T = new(V)
	obj.SetName(nn.Name)
	obj.SetNamespace(nn.Namespace)
	return defaultK.UpdateStatus(obj, func() { f(obj) }, opts...)
}

// Object returns a function that fetches a resource and returns the object.
// It can be used with gomega.Eventually() like this:
//   deployment := appsv1.Deployment{ ... }
//   gomega.Eventually(k.Object(&deployment)).To(HaveField("Spec.Replicas", gomega.Equal(pointer.Int32(3))))
// By calling the returned function directly it can also be used as gomega.Expect(k.Object(...)()).To(...)
func Object[T Obj[V], V any](nn types.NamespacedName) func() (T, error) {
	checkClient()
	return func() (T, error) {
		var obj T = new(V)
		obj.SetName(nn.Name)
		obj.SetNamespace(nn.Namespace)
		err := defaultK.Get(obj)()
		return obj, err
	}
}

// ObjectList returns a function that fetches a resource and returns the object.
// It can be used with gomega.Eventually() like this:
//   deployments := appsv1.DeploymentList{ ... }
//   gomega.Eventually(k.ObjectList(&deployments)).To(HaveField("Items", HaveLen(1)))
// By calling the returned function directly it can also be used as gomega.Expect(k.ObjectList(...)()).To(...)
func ObjectList[T ObjList[V], V client.ObjectList](opts ...client.ListOption) func() (T, error) {
	checkClient()
	return func() (T, error) {
		var obj T = new(V)
		err := defaultK.List(obj, opts...)()
		return obj, err
	}
}
