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
	testiov1 "wanGress/api/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		// 处理删除逻辑，例如清理Envoy配置
		if err := r.cleanupEnvoyConfig(&wangress); err != nil {
			log.Error(err, "清理Envoy配置失败")
			return ctrl.Result{}, err
		}
		// 更新状态并退出
		return ctrl.Result{}, nil
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
	if err := updateWanGressStatus(&wangress); err != nil {
		log.Error(err, "无法更新WanGress状态")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// 生成Envoy配置的函数
func generateEnvoyConfig(wangress *testiov1.WanGress) (EnvoyConfig, error) {
	// 根据wangress对象的规格生成Envoy配置
	// 这里需要根据实际情况实现转换逻辑
	return EnvoyConfig{}, nil
}

// 更新Envoy配置的函数
func updateEnvoyConfig(config EnvoyConfig) error {
	// 实现与Envoy通信，更新配置的逻辑
	// 可能涉及到使用xDS API与Envoy进行通信
	return nil
}

// 更新WanGress状态的函数
func updateEnvoyConfig(config EnvoyConfig) error {
	// 实现与Envoy通信，更新配置的逻辑
	// 可能涉及到使用xDS API与Envoy进行通信
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WanGressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testiov1.WanGress{}).
		Complete(r)
}
