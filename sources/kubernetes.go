package sources

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8sCoreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sSource struct {
	client     k8sCoreV1Types.SecretInterface
	secretName string
}

func newK8sSourceFromURI(uri *url.URL) (Source, error) {
	return newK8sSource(uri.Host, strings.Trim(uri.Path, "/"))
}

func newK8sSource(k8sContext, secretName string) (*k8sSource, error) {
	client, err := newK8sClient(k8sContext)

	if err != nil {
		return nil, err
	}

	return &k8sSource{client, secretName}, nil
}

func (src *k8sSource) ReadSecrets() (secretsMap, error) {
	fmt.Fprintln(os.Stderr, "Reading secrets from Kubernetes…")

	secrets, err := src.client.Get(
		context.Background(),
		src.secretName,
		k8sMetaV1.GetOptions{},
	)

	if err != nil {
		return nil, fmt.Errorf("Error reading secrets from Kubernetes: %w", err)
	}

	return NewSecretsFromByteMap(secrets.Data), nil
}

func (src *k8sSource) WriteSecrets(secrets secretsMap) error {
	fmt.Fprintln(os.Stderr, "Writing secrets to Kubernetes…")

	secretSpec := k8sCoreV1.Secret{
		TypeMeta: k8sMetaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name: src.secretName,
		},
		StringData: map[string]string{},
		Type:       "Opaque",
	}

	_, err := src.client.Update(context.Background(), &secretSpec, k8sMetaV1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("Error writing secrets to Kubernetes: %w", err)
	}

	return nil
}

func newK8sClient(k8sContext string) (k8sCoreV1Types.SecretInterface, error) {
	// It is standard for K8s CLIs to accept a custom KUBECONFIG path.
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		// The default location for K8s config.
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: k8sContext,
			// ClusterInfo: clientcmdapi.Cluster{Server: ""},
		}).ClientConfig()

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, fmt.Errorf("Error creating Kubernetes client: %w", err)
	}

	return clientset.CoreV1().Secrets("default"), nil
}
