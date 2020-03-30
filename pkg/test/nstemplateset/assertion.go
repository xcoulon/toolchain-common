package nstemplateset

import (
	"context"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/pkg/apis/toolchain/v1alpha1"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Assertion struct {
	nsTmplSet      *toolchainv1alpha1.NSTemplateSet
	client         client.Client
	namespacedName types.NamespacedName
	t              test.T
}

func (a *Assertion) loadNSTemplateSet() error {
	if a.nsTmplSet != nil {
		return nil
	}
	nsTmplSet := &toolchainv1alpha1.NSTemplateSet{}
	err := a.client.Get(context.TODO(), a.namespacedName, nsTmplSet)
	a.nsTmplSet = nsTmplSet
	return err
}

func AssertThatNSTemplateSet(t test.T, namespace, name string, client client.Client) *Assertion {
	return &Assertion{
		client:         client,
		namespacedName: test.NamespacedName(namespace, name),
		t:              t,
	}
}

func (a *Assertion) HasNoConditions() *Assertion {
	err := a.loadNSTemplateSet()
	require.NoError(a.t, err)
	require.Empty(a.t, a.nsTmplSet.Status.Conditions)
	return a
}

func (a *Assertion) HasConditions(expected ...toolchainv1alpha1.Condition) *Assertion {
	err := a.loadNSTemplateSet()
	require.NoError(a.t, err)
	test.AssertConditionsMatch(a.t, a.nsTmplSet.Status.Conditions, expected...)
	return a
}

func Provisioned() toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:   toolchainv1alpha1.ConditionReady,
		Status: corev1.ConditionTrue,
		Reason: toolchainv1alpha1.NSTemplateSetProvisionedReason,
	}
}

func Provisioning() toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:   toolchainv1alpha1.ConditionReady,
		Status: corev1.ConditionFalse,
		Reason: toolchainv1alpha1.NSTemplateSetProvisioningReason,
	}
}

func Updating() toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:   toolchainv1alpha1.ConditionReady,
		Status: corev1.ConditionFalse,
		Reason: toolchainv1alpha1.NSTemplateSetUpdatingReason,
	}
}

func UpdateFailed(msg string) toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:    toolchainv1alpha1.ConditionReady,
		Status:  corev1.ConditionFalse,
		Reason:  toolchainv1alpha1.NSTemplateSetUpdateFailedReason,
		Message: msg,
	}
}

func UnableToProvision(msg string) toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:    toolchainv1alpha1.ConditionReady,
		Status:  corev1.ConditionFalse,
		Reason:  toolchainv1alpha1.NSTemplateSetUnableToProvisionReason,
		Message: msg,
	}
}

func UnableToProvisionNamespace(msg string) toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:    toolchainv1alpha1.ConditionReady,
		Status:  corev1.ConditionFalse,
		Reason:  toolchainv1alpha1.NSTemplateSetUnableToProvisionNamespaceReason,
		Message: msg,
	}
}

func Terminating() toolchainv1alpha1.Condition {
	return toolchainv1alpha1.Condition{
		Type:   toolchainv1alpha1.ConditionReady,
		Status: corev1.ConditionFalse,
		Reason: toolchainv1alpha1.NSTemplateSetTerminatingReason,
	}
}

func (a *Assertion) HasFinalizer() *Assertion {
	err := a.loadNSTemplateSet()
	require.NoError(a.t, err)
	assert.Len(a.t, a.nsTmplSet.Finalizers, 1)
	assert.Contains(a.t, a.nsTmplSet.Finalizers, "finalizer.toolchain.dev.openshift.com")
	return a
}

func (a *Assertion) DoesNotHaveFinalizer() *Assertion {
	err := a.loadNSTemplateSet()
	require.NoError(a.t, err)
	assert.Len(a.t, a.nsTmplSet.Finalizers, 0)
	return a
}