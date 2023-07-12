package server

import (
	"context"
	"io/ioutil"
	"net/http"
	types "random_name_controller/pkg/apis/v1"
	"time"

	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

var (
	defaultRules = types.RandomStringSetEntryRules{
		Length:      20,
		Digits:      10,
		Symbols:     0,
		AllowUpper:  false,
		AllowRepeat: false,
	}
)

type SyncRequest struct {
	Parent   types.RandomStringSet `json:"parent"`
	Children SyncRequestChildren   `json:"children"`
}

type SyncRequestChildren struct {
	ConfigMaps map[string]*v1.ConfigMap `json:"ConfigMap.v1"`
}

type SyncResponse struct {
	Status   *types.RandomStringSetStatus `json:"status"`
	Children []runtime.Object             `json:"children"`
}

type HostArgs struct {
	Logger        *zap.Logger
	ListenAddress string
}

type PodIdentityArgs struct {
	UsePodIdentity bool
	Label          string
	Value          string
	UseVmIdentity  bool
}
type Host interface {
	Run(ctx context.Context) error
}

type host struct {
	cfg *HostArgs
}

func New(in *HostArgs) Host {
	return &host{
		cfg: in,
	}
}

func (h *host) Run(ctx context.Context) error {
	errorChannel := make(chan error)
	http.HandleFunc("/sync", h.SyncHandler)
	go func() {
		err := http.ListenAndServe(h.cfg.ListenAddress, nil)
		errorChannel <- err
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errorChannel:
			return err
		}
	}
}

func (h *host) sync(ctx context.Context, request *SyncRequest) (*SyncResponse, error) {
	response := &SyncResponse{}

	cfg, status := GenerateConfigMap(&request.Parent)
	if cfg != nil {
		response.Children = []runtime.Object{cfg}
	}
	response.Status = &status
	return response, nil
}

func (h *host) SyncHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request := &SyncRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.sync(req.Context(), request)
	if err != nil {
		h.cfg.Logger.Error("failed to synchronize results", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func GenerateConfigMap(source *types.RandomStringSet) (*v1.ConfigMap, types.RandomStringSetStatus) {
	cfg := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{},
			Name:      source.Spec.Name,
			Namespace: source.Namespace,
		},
		Data: map[string]string{},
	}

	for _, entry := range source.Spec.Entries {
		var rules *types.RandomStringSetEntryRules
		if entry.Rules == nil {
			rules = &defaultRules
		} else {
			rules = entry.Rules
		}

		value, err := password.Generate(rules.Length, rules.Digits, rules.Symbols, rules.AllowUpper, rules.AllowRepeat)
		if err != nil {
			return nil, types.RandomStringSetStatus{
				Conditions: []metav1.Condition{
					{
						Type:               "Ready",
						Status:             metav1.ConditionFalse,
						ObservedGeneration: source.ObjectMeta.Generation,
						LastTransitionTime: metav1.Time{
							Time: time.Now(),
						},
						Reason:  "Error",
						Message: err.Error(),
					},
				},
				LastObservedGeneration: source.ObjectMeta.Generation,
			}
		}

		cfg.Data[entry.Name] = value
	}

	return cfg, types.RandomStringSetStatus{
		Conditions: []metav1.Condition{
			{
				Type:               "Ready",
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 0,
				LastTransitionTime: metav1.Time{
					Time: time.Now(),
				},
				Reason:  "Success",
				Message: "Success",
			},
		},
		LastObservedGeneration: source.ObjectMeta.Generation,
	}
}
