package k8s

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	*kubernetes.Clientset
}

var k8sClientInstance *K8sClient
var once sync.Once

func getKubeConfig() (*rest.Config, error) {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func GetK8SClient() (*K8sClient, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	once.Do(func() {
		k8sClientInstance = &K8sClient{clientSet}
	})
	return k8sClientInstance, err
}

var (
	INGRESS_CLASS_NAME = "nginx"
)

func (k *K8sClient) CreateNamespaceIfNotExists(namespaceName string) error {
	_, err := k.CoreV1().Namespaces().Get(context.Background(), namespaceName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		namespace := &apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		_, err := k.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
		return err
	}
	return nil
}

func (k *K8sClient) CreateDeployment(name, namespace string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            name,
							Image:           fmt.Sprintf("docker.io/library/%s", name),
							ImagePullPolicy: apiv1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	_, err := k.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	return err
}

func (k *K8sClient) ExposeApp(name, namespace, domain string) error {

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-svc", name),
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:     "http",
					Port:     3000,
					Protocol: apiv1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "http",
					},
				},
			},
			Selector: map[string]string{
				"app": name,
			},
		},
	}
	_, err := k.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	pathType := networkv1.PathTypePrefix
	ingress := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-ingress", name),
			Namespace: namespace,
			Annotations: map[string]string{
				"cert-manager.io/cluster-issuer": "letsencrypt-issuer",
			},
		},
		Spec: networkv1.IngressSpec{
			IngressClassName: &INGRESS_CLASS_NAME,
			Rules: []networkv1.IngressRule{
				{
					Host: fmt.Sprintf("%s.%s", name, domain),
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkv1.IngressBackend{
										Service: &networkv1.IngressServiceBackend{
											Name: service.Name,
											Port: networkv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = k.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
	return err
}
