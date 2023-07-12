package server_test

import (
	types "random_name_controller/pkg/apis/v1"
	"random_name_controller/pkg/server"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidDoc(t *testing.T) {
	configMapName := strconv.Itoa(int(time.Now().UnixMilli()))
	input := &types.RandomStringSet{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "testing"},
		Spec: types.RandomStringSetSpec{Name: configMapName, Entries: []types.RandomStringSetEntry{{Name: "first", Rules: &types.RandomStringSetEntryRules{
			Length:      30,
			Digits:      10,
			Symbols:     5,
			AllowUpper:  false,
			AllowRepeat: false,
		},
		},
			{Name: "second"}}},
		Status: types.RandomStringSetStatus{},
	}
	cfg, status := server.GenerateConfigMap(input)
	assert.NotNil(t, cfg, "configmap should not be nil")
	assert.Equal(t, input.Namespace, cfg.Namespace, "namespaces should match")
	assert.Equal(t, input.Spec.Name, cfg.Name, "name should match specified name")
	assert.Equal(t, 1, len(status.Conditions))
	assert.Equal(t, metav1.ConditionTrue, status.Conditions[0].Status)
}
