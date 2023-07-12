package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RandomStringSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   RandomStringSetSpec   `json:"spec"`
	Status RandomStringSetStatus `json:"status,omitempty"`
}

type RandomStringSetSpec struct {
	Name    string                 `json:"name"`
	Entries []RandomStringSetEntry `json:"entries"`
}

type RandomStringSetEntry struct {
	Name  string                     `json:"name"`
	Rules *RandomStringSetEntryRules `json:"rules"`
}

type RandomStringSetEntryRules struct {
	Length      int  `json:"length"`
	Digits      int  `json:"digits"`
	Symbols     int  `json:"symbols"`
	AllowUpper  bool `json:"allowUpper"`
	AllowRepeat bool `json:"allowRepeat"`
}

type RandomStringSetStatus struct {
	Conditions             []metav1.Condition `json:"conditions,omitempty"`
	LastObservedGeneration int64              `json:"lastObservedGeneration"`
}
