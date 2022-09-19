/*
Copyright 2021.

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

package controllers

import (
	"context"
	"runtime"

	mpsv1alpha1 "github.com/playfab/thundernetes/pkg/operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type DynamicStandbyReconciler struct {
	client.Client
	Scheme   *k8sruntime.Scheme
	Recorder record.EventRecorder
}

func NewDynamicStandbyReconciler(mgr manager.Manager) *DynamicStandbyReconciler {
	cl := mgr.GetClient()
	return &DynamicStandbyReconciler{
		Client:   cl,
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("DynamicStandby"),
	}
}

func (r *DynamicStandbyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var gsb mpsv1alpha1.GameServerBuild
	if err := r.Get(ctx, req.NamespacedName, &gsb); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to fetch GameServerBuild - it is being deleted")
			// GameServerBuild is being deleted so clear its entry from the crashesPerBuild map
			// no-op if the entry is not present
			crashesPerBuild.Delete(getKeyForCrashesPerBuildMap(&gsb))
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch gameServerBuild")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *DynamicStandbyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mpsv1alpha1.GameServerBuild{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: runtime.NumCPU(),
		}).
		Complete(r)
}
