package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"github.com/kyverno/chainsaw/pkg/utils/kube"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/pkg/ext/resource/convert"
	corev1 "k8s.io/api/core/v1"
)

func buildNamespace(ctx context.Context, compilers compilers.Compilers, name string, template *v1alpha1.Projection, tc apis.Bindings) (*corev1.Namespace, error) {
	namespace := kube.Namespace(name)
	if template == nil {
		return &namespace, nil
	}
	if template.Value() == nil {
		return &namespace, nil
	}
	object := kube.ToUnstructured(&namespace)
	tc = bindings.RegisterBinding(ctx, tc, "namespace", object.GetName())
	merged, err := templating.TemplateAndMerge(ctx, compilers, object, tc, *template)
	if err != nil {
		return nil, err
	}
	return convert.To[corev1.Namespace](merged)
}
