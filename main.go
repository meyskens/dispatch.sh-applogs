package main

import (
	"os"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func main() {
	client, err := newKubernetesClient()
	if err != nil {
		panic(err)
	}

	watchlist := cache.NewListWatchFromClient(client.Core().RESTClient(), "pods", os.Getenv("MY_NAMESPACE"), fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if pod.Status.Phase == "Running" {
					for _, container := range pod.Spec.Containers {
						go follow(pod.GetObjectMeta().GetLabels()["release"], pod.GetName(), container.Name)
					}
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod := newObj.(*v1.Pod)
				if pod.Status.Phase == "Running" {
					for _, container := range pod.Spec.Containers {
						go follow(pod.GetObjectMeta().GetLabels()["release"], pod.GetName(), container.Name)
					}
				}
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	<-stop
}
