package seleniumhub

import (
	"context"

	testv1alpha1 "github.com/WianVos/selenium_k8s_operator/pkg/apis/test/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_seleniumhub")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SeleniumHub Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSeleniumHub{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("seleniumhub-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SeleniumHub
	err = c.Watch(&source.Kind{Type: &testv1alpha1.SeleniumHub{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SeleniumHub
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &testv1alpha1.SeleniumHub{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSeleniumHub implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSeleniumHub{}

// ReconcileSeleniumHub reconciles a SeleniumHub object
type ReconcileSeleniumHub struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SeleniumHub object and makes changes based on the state read
// and what is in the SeleniumHub.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSeleniumHub) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SeleniumHub")

	// Fetch the SeleniumHub instance
	instance := &testv1alpha1.SeleniumHub{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	selHub := instance
	podList := &corev1.PodList{}
	reqLogger.Info()
	lbs := map[string]string{
		"app": instance.Name,
	}

	labelSelector := labels.SelectorFromSet(lbs)
	listOps := &client.ListOptions{Namespace: selHub.Namespace, LabelSelector: labelSelector}
	if err = r.client.List(context.TODO(), podList, listOps); err != nil {
		return reconcile.Result{}, err
	}

	for _, pod := range podList.Items {
		reqLogger.Info("found this pod!!!!", "pod", pod)
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set SeleniumHub instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check if this pod is up to spec

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *testv1alpha1.SeleniumHub) *corev1.Pod {
	labels := map[string]string{
		"app":  cr.Name,
		"role": "hub",
	}

	return &corev1.Pod{

		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "selenium-hub",
					Image: "selenium/hub:3.141",
					Ports: []corev1.ContainerPort{{
						ContainerPort: 4444,
						Name:          "selenium",
					}},
					Resources: getResourceRequirements(getResourceList(cr.Spec.CPU, cr.Spec.Memory), getResourceList("", "")),
				},
			},
		},
	}
}

func getResourceList(cpu, memory string) v1.ResourceList {
	res := v1.ResourceList{}
	if cpu != "" {
		res[v1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[v1.ResourceMemory] = resource.MustParse(memory)
	}
	return res
}

func getResourceRequirements(requests, limits v1.ResourceList) v1.ResourceRequirements {
	res := v1.ResourceRequirements{}
	res.Requests = requests
	res.Limits = limits
	return res
}

// apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: selenium-hub
//   labels:
//     app: selenium-hub
// spec:
//   replicas: 1
//   selector:
//     matchLabels:
//       app: selenium-hub
//   template:
//     metadata:
//       labels:
//         app: selenium-hub
//     spec:
//       containers:
//       - name: selenium-hub
//         image: selenium/hub:3.141
//         ports:
//           - containerPort: 4444
//         resources:
//           limits:
//             memory: "1000Mi"
//             cpu: ".5"
//         livenessProbe:
//           httpGet:
//             path: /wd/hub/status
//             port: 4444
//           initialDelaySeconds: 30
//           timeoutSeconds: 5
//         readinessProbe:
//           httpGet:
//             path: /wd/hub/status
//             port: 4444
//           initialDelaySeconds: 30
//           timeoutSeconds: 5
