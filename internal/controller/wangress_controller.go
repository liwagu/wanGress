package controller

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	testiov1 "wanGress/api/v1"
	"wanGress/pkg/envoy"
)

type WanGressReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	XDSClient *envoy.XDSClient
}

func (r *WanGressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("wangress", req.NamespacedName)

	var wangress testiov1.WanGress
	if err := r.Get(ctx, req.NamespacedName, &wangress); err != nil {
		log.Error(err, "unable to fetch WanGress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if wangress.DeletionTimestamp != nil {
		if err := r.cleanupEnvoyConfig(ctx, &wangress); err != nil {
			log.Error(err, "failed to cleanup Envoy configuration")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	listeners, clusters, routes, endpoints, err := envoy.GenerateEnvoyConfig(&wangress)
	if err != nil {
		log.Error(err, "failed to generate Envoy configuration")
		return ctrl.Result{}, err
	}

	if err := r.XDSClient.UpdateConfig(listeners, clusters, routes, endpoints); err != nil {
		log.Error(err, "failed to update Envoy configuration")
		return ctrl.Result{}, err
	}

	wangress.Status.Phase = "Active"
	wangress.Status.Message = "Envoy configuration applied successfully"
	if err := r.Status().Update(ctx, &wangress); err != nil {
		log.Error(err, "failed to update WanGress status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WanGressReconciler) cleanupEnvoyConfig(ctx context.Context, wangress *testiov1.WanGress) error {
	return r.XDSClient.RemoveConfig()
}

func (r *WanGressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.XDSClient = envoy.NewXDSClient()
	r.XDSClient.Run()

	return ctrl.NewControllerManagedBy(mgr).
		For(&testiov1.WanGress{}).
		Complete(r)
}
