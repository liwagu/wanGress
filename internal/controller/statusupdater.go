package controller

import (
	"context"
	testiov1 "wanGress/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func updateWanGressStatus(client client.Client, wangress *testiov1.WanGress, phase string, message string, condition metav1.Condition) error {
	// Update Phase and Message directly
	wangress.Status.Phase = phase
	wangress.Status.Message = message

	// Example of updating Conditions
	// This is a simplified example. In practice, you would likely search for an existing condition of the same type and update it
	wangress.Status.Conditions = append(wangress.Status.Conditions, condition)

	return client.Status().Update(context.Background(), wangress)
}
