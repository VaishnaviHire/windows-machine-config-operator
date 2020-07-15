package e2e

import (
	"context"
	"testing"

	mapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func deletionTestSuite(t *testing.T) {
	t.Run("Deletion", func(t *testing.T) { testWindowsNodeDeletion(t) })
}

// testWindowsNodeDeletion tests the Windows node deletion from the cluster.
func testWindowsNodeDeletion(t *testing.T) {
	testCtx, err := NewTestContext(t)
	require.NoError(t, err)

	// The second item in the slice is the Windows MachineSet we're interested in scaling down.
	require.Len(t, gc.machineSetNames, 2)
	windowsMachineSet := &mapi.MachineSet{}
	err = framework.Global.Client.Get(context.TODO(), types.NamespacedName{Name: gc.machineSetNames[1],
		Namespace: "openshift-machine-api"}, windowsMachineSet)
	// Reset the number of nodes to be deleted to 0
	gc.numberOfNodes = 0
	// Delete the Windows VM that got created.
	windowsMachineSet.Spec.Replicas = &gc.numberOfNodes
	if err := framework.Global.Client.Update(context.TODO(), windowsMachineSet); err != nil {
		t.Fatalf("error updating windowsMachineSet custom resource  %v", err)
	}
	// As per testing, each windows VM is taking roughly 12 minutes to be shown up in the cluster
	err = testCtx.waitForWindowsNodes(gc.numberOfNodes, true, true)
	if err != nil {
		t.Fatalf("windows node deletion failed  with %v", err)
	}
	// Clean up the windows MachineSets created by the test suite
	for _, machineSetName := range gc.machineSetNames {
		machineSet := &mapi.MachineSet{}
		err = framework.Global.Client.Get(context.TODO(), types.NamespacedName{Name: machineSetName,
			Namespace: "openshift-machine-api"}, machineSet)
		if err := framework.Global.Client.Delete(context.TODO(), windowsMachineSet); err != nil {
			t.Fatalf("Windows MachineSet")
		}
	}
}
