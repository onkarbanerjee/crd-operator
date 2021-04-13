package handler

import (
	"os"

	v1 "github.com/onkarbanerjee/crd-operator/pkg/apis/customconfig/v1"
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
	ObjectUpdated(obj interface{})
}

// CCHandler is a sample implementation of Handler
type CCHandler struct {
	Client kubernetes.Interface
}

// Init handles any handler initialization
func (t *CCHandler) Init() error {
	log.Info("TestHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *CCHandler) ObjectCreated(obj interface{}) {
	log.Info("CCHandler.ObjectCreated")

	cc, ok1 := obj.(*v1.CustomConfig)
	log.Info("ok is", ok1)

	log.Info("cc is ", cc.Spec.Key, cc.Spec.Value, cc.Spec.ConfigmapName)

	ns := os.Getenv("NAMESPACE")
	log.Info("ns is", ns)
	if ns == "" {
		ns = "default"
	}

	cm := core_v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      cc.Spec.ConfigmapName,
			Namespace: ns,
		},
		Data: map[string]string{
			cc.Spec.Key: cc.Spec.Value,
		},
	}

	// if t.Client == nil {
	// 	log.Info("handler doesnt have a k8s client, returning")
	// 	return
	// }
	_, err := t.Client.CoreV1().ConfigMaps(ns).Create(&cm)
	if err != nil {
		if errors.IsAlreadyExists(err) || errors.IsConflict(err) {
			// Resource already exists. Carry on.
			log.Info("config map is already exist")
			return
		}
		log.Error("error is", err)
		return
	}

	log.Info("config map created")
}

// ObjectDeleted is called when an object is deleted
func (t *CCHandler) ObjectDeleted(obj interface{}) {
	log.Info("CCHandler.ObjectDeleted")
	cc, ok1 := obj.(*v1.CustomConfig)
	log.Info("ok is", ok1)

	log.Info("cc is ", cc.Spec.Key, cc.Spec.Value, cc.Spec.ConfigmapName)

	ns := os.Getenv("NAMESPACE")
	log.Info("ns is", ns)
	if ns == "" {
		ns = "default"
	}

	cm, err := t.Client.CoreV1().ConfigMaps(ns).Get(cc.Spec.ConfigmapName, meta_v1.GetOptions{})
	if err != nil {
		log.Error("error is", err)
		return
	}

	delete(cm.Data, cc.Spec.Key)

	log.Info("len of data is", len(cm.Data))
	if len(cm.Data) == 0 {
		err = t.Client.CoreV1().ConfigMaps(ns).Delete(cm.Name, nil)
		if err != nil {
			log.Error("error is", err)
		}
		return
	}

	_, err = t.Client.CoreV1().ConfigMaps(ns).Update(cm)

	if err != nil {
		log.Error("error is", err)
		return
	}

}

// ObjectUpdated is called when an object is updated
func (t *CCHandler) ObjectUpdated(obj interface{}) {
	log.Info("CCHandler.ObjectUpdated")

	cc, ok1 := obj.(*v1.CustomConfig)
	log.Info("ok is", ok1)

	log.Info("cc is ", cc.Spec.Key, cc.Spec.Value, cc.Spec.ConfigmapName)

	ns := os.Getenv("NAMESPACE")
	log.Info("ns is", ns)
	if ns == "" {
		ns = "default"
	}

	cm, err := t.Client.CoreV1().ConfigMaps(ns).Get(cc.Spec.ConfigmapName, meta_v1.GetOptions{})
	if err != nil {
		log.Error("error is", err)
		return
	}

	cm.Data[cc.Spec.Key] = cc.Spec.Value
	_, err = t.Client.CoreV1().ConfigMaps(ns).Update(cm)

	if err != nil {
		log.Error("error is", err)
		return
	}
}
