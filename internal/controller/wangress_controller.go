/*
Copyright 2024.

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

package controller

import (
	"context"
	envoy_api_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/go-logr/logr"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	testiov1 "wanGress/api/v1"
)

// WanGressReconciler reconciles a WanGress object
type WanGressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.io.liwa.com,resources=wangresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.io.liwa.com,resources=wangresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.io.liwa.com,resources=wangresses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WanGress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *WanGressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("wangress", req.NamespacedName)

	// 获取WanGress实例
	var wangress testiov1.WanGress
	if err := r.Get(ctx, req.NamespacedName, &wangress); err != nil {
		log.Error(err, "无法获取WanGress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 检查WanGress是否被标记为删除
	if wangress.DeletionTimestamp != nil {
		if err := r.cleanupEnvoyConfig(ctx, &wangress, log); err != nil {
			log.Error(err, "Failed to cleanup Envoy configuration")
			return ctrl.Result{}, err
		}
		// 其他删除逻辑...
	}

	// 基于WanGress规格生成Envoy配置
	envoyConfig, err := generateEnvoyConfig(&wangress)
	if err != nil {
		log.Error(err, "无法生成Envoy配置")
		return ctrl.Result{}, err
	}

	// 更新Envoy配置
	if err := updateEnvoyConfig(envoyConfig); err != nil {
		log.Error(err, "无法更新Envoy配置")
		return ctrl.Result{}, err
	}

	// 更新WanGress状态
	if err := updateWanGressStatus(r.Client, &wangress, "Active", "Envoy configuration applied successfully.", condition); err != nil {
		log.Error(err, "Failed to update WanGress status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// 生成Envoy配置的函数
func generateEnvoyConfig(wangress *testiov1.WanGress) (*envoy_listener_v3.Listener, error) {
	// 根据WanGress规格生成Envoy配置
	// 这里只是一个示例，您需要根据实际需求来实现
	manager := &envoy_http_connection_manager_v3.HttpConnectionManager{
		CodecType:  envoy_http_connection_manager_v3.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http",
		// 添加更多的配置逻辑
	}

	pbst, err := anypb.New(manager)
	if err != nil {
		return nil, err
	}

	listener := &envoy_listener_v3.Listener{
		Name: "listener_0",
		Address: &envoy_api_v3.Address{
			Address: &envoy_api_v3.Address_SocketAddress{
				SocketAddress: &envoy_api_v3.SocketAddress{
					Protocol: envoy_api_v3.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &envoy_api_v3.SocketAddress_PortValue{
						PortValue: 80,
					},
				},
			},
		},
		FilterChains: []*envoy_listener_v3.FilterChain{
			{
				Filters: []*envoy_listener_v3.Filter{
					{
						Name: "envoy.filters.network.http_connection_manager",
						ConfigType: &envoy_listener_v3.Filter_TypedConfig{
							TypedConfig: pbst,
						},
					},
				},
			},
		},
	}

	return listener, nil
}

// 更新Envoy配置的函数
func updateEnvoyConfig(config EnvoyConfig) error {
	// 实现与Envoy通信，更新配置的逻辑
	// 可能涉及到使用xDS API与Envoy进行通信
	return nil
}


// cleanupEnvoyConfig cleans up Envoy configuration for a given WanGress resource.
// This is a placeholder function and needs to be implemented based on your Envoy configuration strategy.
func (r *WanGressReconciler) cleanupEnvoyConfig(ctx context.Context, wangress *testiov1.WanGress, log logr.Logger) error {
	// Log the cleanup operation
	log.Info("Cleaning up Envoy configuration", "WanGress", wangress.Name)

	// Implement the logic to communicate with Envoy or its management service to remove the configuration
	// associated with the WanGress resource. This might involve calling Envoy's API, modifying config files,
	// or interacting with a service like Contour that manages Envoy configuration.

	// Placeholder for cleanup logic
	// err := removeEnvoyConfiguration(wangress)
	// if err != nil {
	//     log.Error(err, "Failed to remove Envoy configuration", "WanGress", wangress.Name)
	//     return err
	// }

	// Log the successful cleanup
	log.Info("Successfully cleaned up Envoy configuration", "WanGress", wangress.Name)

	return nil
}


// SetupWithManager 注册WanGress资源，使控制器能够监视资源的变化
func (r *WanGressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testiov1.WanGress{}).
		Complete(r)
}
