/*
Copyright Yann Lacroix.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	context "context"
	time "time"

	apiskubeconfigv1 "github.com/brawdunoir/kubebrowser/pkg/apis/kubeconfig/v1"
	versioned "github.com/brawdunoir/kubebrowser/pkg/client/clientset/versioned"
	internalinterfaces "github.com/brawdunoir/kubebrowser/pkg/client/informers/externalversions/internalinterfaces"
	kubeconfigv1 "github.com/brawdunoir/kubebrowser/pkg/client/listers/kubeconfig/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// KubeconfigInformer provides access to a shared informer and lister for
// Kubeconfigs.
type KubeconfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() kubeconfigv1.KubeconfigLister
}

type kubeconfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewKubeconfigInformer constructs a new informer for Kubeconfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKubeconfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKubeconfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredKubeconfigInformer constructs a new informer for Kubeconfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKubeconfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeconfigV1().Kubeconfigs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeconfigV1().Kubeconfigs(namespace).Watch(context.TODO(), options)
			},
		},
		&apiskubeconfigv1.Kubeconfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *kubeconfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKubeconfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *kubeconfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiskubeconfigv1.Kubeconfig{}, f.defaultInformer)
}

func (f *kubeconfigInformer) Lister() kubeconfigv1.KubeconfigLister {
	return kubeconfigv1.NewKubeconfigLister(f.Informer().GetIndexer())
}
