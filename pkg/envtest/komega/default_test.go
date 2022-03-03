package komega

import (
	"testing"

	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

func TestDefaultObject(t *testing.T) {
	g := NewWithT(t)

	fc := createFakeClient()
	SetClient(fc)

	g.Eventually(Object[*appsv1.Deployment](types.NamespacedName{Namespace: "default", Name: "test"})).Should(And(
		Not(BeNil()),
		HaveField("Spec.Replicas", Equal(pointer.Int32(5))),
	))
}
