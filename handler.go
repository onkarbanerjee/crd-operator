package main

import (
	v1 "github.com/onkarbanerjee/crd-custom-config/pkg/apis/customconfig/v1"
	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

// TestHandler is a sample implementation of Handler
type TestHandler struct {
	Client kubernetes.Interface
}

// Init handles any handler initialization
func (t *TestHandler) Init() error {
	log.Info("TestHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *TestHandler) ObjectCreated(obj interface{}) {
	log.Info("TestHandler.ObjectCreated")

	cc, ok1 := obj.(*v1.CustomConfig)
	log.Info("ok is", ok1)

	log.Info("cc is ", cc.Spec.Key, cc.Spec.Value, cc.Spec.ConfigmapName)

	cm := core_v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      cc.Spec.ConfigmapName,
			Namespace: "default",
			// Labels:    map[string]string{}"custom-config-cm",
		},
		Data: map[string]string{
			// "__CFG_TYPE":        "mgmt-cfg",
			// "cim.json":          string(cimData),
			// "cim.yang":          string(yangData),
			// "updatePolicy.json": string(cimUpdatePolicyJson),
			cc.Spec.Key: cc.Spec.Value,
			"revision":  "0",
		},
	}

	// if t.Client == nil {
	// 	log.Info("handler doesnt have a k8s client, returning")
	// 	return
	// }
	_, err := t.Client.CoreV1().ConfigMaps("default").Create(&cm)
	if errors.IsAlreadyExists(err) || errors.IsConflict(err) {
		// Resource already exists. Carry on.
		log.Info("config map is already exist")
		return
	}
	log.Info("config map created")
}

// ObjectDeleted is called when an object is deleted
func (t *TestHandler) ObjectDeleted(obj interface{}) {
	log.Info("TestHandler.ObjectDeleted")
}

// ObjectUpdated is called when an object is updated
func (t *TestHandler) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("TestHandler.ObjectUpdated")
}
