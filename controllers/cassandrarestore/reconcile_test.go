package cassandrarestore

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cscetbon/casskop/controllers/common"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"

	api "github.com/cscetbon/casskop/api/v2"
	"github.com/ghodss/yaml"
)

var cassandraRestoreYaml = `
apiVersion: db.orange.com/v2
kind: CassandraRestore
metadata:
  name: test-cassandra-restore
spec:
  cassandraCluster: test-cluster-dc1
  cassandraBackup: test-cassandra-backup
  concurrentConnection: 15
  noDeleteTruncates: false
#   schemaVersion:
#   exactSchemaVersion:
  entities: "k1,k2.t1"
`

func helperInitCassandraRestore(cassandraRestoreYaml string) api.CassandraRestore {
	var cassandraRestore api.CassandraRestore
	if err := yaml.Unmarshal([]byte(cassandraRestoreYaml), &cassandraRestore); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	return cassandraRestore
}

func helperInitCassandraRestoreController(cassandraRestoreYaml string) (*CassandraRestoreReconciler,
	*api.CassandraRestore, *record.FakeRecorder) {
	//cassandraBackup := common.HelperInitCassandraBackup(cassandraRestoreYaml)
	cassandraRestore := helperInitCassandraRestore(cassandraRestoreYaml)

	cassandraRestoreList := api.CassandraRestoreList{}

	// Register operator types with the runtime scheme.
	fakeClientScheme := scheme.Scheme
	fakeClientScheme.AddKnownTypes(api.GroupVersion, &api.CassandraCluster{})
	fakeClientScheme.AddKnownTypes(api.GroupVersion, &api.CassandraClusterList{})
	fakeClientScheme.AddKnownTypes(api.GroupVersion, &api.CassandraBackup{})
	fakeClientScheme.AddKnownTypes(api.GroupVersion, &api.CassandraBackupList{})
	fakeClientScheme.AddKnownTypes(api.GroupVersion, &cassandraRestore)
	fakeClientScheme.AddKnownTypes(api.GroupVersion, &cassandraRestoreList)

	objs := []runtime.Object{
		&cassandraRestore,
	}

	fakeClient := fake.NewClientBuilder().WithScheme(fakeClientScheme).WithRuntimeObjects(objs...).Build()

	fakeRecorder := record.NewFakeRecorder(3)
	CassandraRestoreReconciler := CassandraRestoreReconciler{
		Client:   fakeClient,
		Scheme:   fakeClientScheme,
		Recorder: fakeRecorder,
	}

	return &CassandraRestoreReconciler, &cassandraRestore, fakeRecorder
}

func TestCassandraRestoreWithUnknownCassandraCluster(t *testing.T) {
	assert := assert.New(t)
	CassandraRestoreReconciler, cassandraRestore, recorder := helperInitCassandraRestoreController(cassandraRestoreYaml)

	CassandraRestoreReconciler.Client.Create(context.TODO(), cassandraRestore)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      cassandraRestore.Name,
			Namespace: cassandraRestore.Namespace,
		},
	}

	res, err := CassandraRestoreReconciler.Reconcile(context.TODO(), req)

	assert.Equal(reconcile.Result{}, res)
	assert.NotNil(err)
	assert.Equal(err.Error(), fmt.Sprintf("cassandraclusters.db.orange.com \"%s\" not found",
		cassandraRestore.Spec.CassandraCluster))
	common.AssertEvent(t, recorder.Events,
		fmt.Sprintf("Warning CassandraClusterNotFound Cassandra Cluster %s to restore not found",
			cassandraRestore.Spec.CassandraCluster))
}

func TestCassandraRestoreWithUnknownCassandraBackup(t *testing.T) {
	assert := assert.New(t)

	CassandraRestoreReconciler, cassandraRestore, recorder := helperInitCassandraRestoreController(cassandraRestoreYaml)

	cassandraCluster := api.CassandraCluster{}
	cassandraCluster.Name = cassandraRestore.Spec.CassandraCluster
	cassandraCluster.Namespace = cassandraRestore.Namespace
	fmt.Println(cassandraCluster)

	CassandraRestoreReconciler.Client.Create(context.TODO(), &cassandraCluster)
	CassandraRestoreReconciler.Client.Create(context.TODO(), cassandraRestore)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      cassandraRestore.Name,
			Namespace: cassandraRestore.Namespace,
		},
	}

	res, err := CassandraRestoreReconciler.Reconcile(context.TODO(), req)

	assert.Equal(reconcile.Result{}, res)
	assert.NotNil(err)
	assert.Equal(fmt.Sprintf("cassandrabackups.db.orange.com \"%s\" not found",
		cassandraRestore.Spec.CassandraBackup), err.Error())
	common.AssertEvent(t, recorder.Events,
		fmt.Sprintf("Warning BackupNotFound Backup %s to restore not found",
			cassandraRestore.Spec.CassandraBackup))
}

func TestCassandraRestoreWithNilStatusCondition(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}

func TestCassandraRestoreWithNoCoordinatorMember(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}

func TestCassandraRestorePhaseRequiredButNoPods(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}

func TestCassandraRestorePhaseRequired(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}

func TestCassandraRestorePhaseInProgress(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}
