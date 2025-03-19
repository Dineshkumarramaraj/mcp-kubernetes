package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// wanted to separate the Kubernetes components from the MCP server.
type kubeClient struct {
	Path   string
	Config *rest.Config
}

var kubeConfig *kubeClient

func (k *kubeClient) dynamicClient() (*dynamic.DynamicClient, error) {
	dynamicClient, err := dynamic.NewForConfig(k.Config)
	if err != nil {
		return nil, err
	}

	return dynamicClient, nil
}

func (k *kubeClient) getResource(gvk *schema.GroupVersionKind) (*schema.GroupVersionResource, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(k.Config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
	m, err := mapper.RESTMapping(schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}, gvk.Version)
	if err != nil {
		return nil, err
	}

	return &m.Resource, nil
}

func NewKubernetesClient() (*kubeClient, error) {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return &kubeClient{
		Path:   kubeconfigPath,
		Config: config,
	}, nil
}

func init() {
	//wanted to separate the Kubernetes components from the MCP server. So,
	//getting kubeconfig file before executing the mcp server
	var err error
	kubeConfig, err = NewKubernetesClient()
	if err != nil {
		log.Fatalf("Error creating Kubernetes Client: %v", err)
	}
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-kubernetes-listpods",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("Kubernetes_listpods",
		mcp.WithDescription("List the pods available in given namespace"),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description("Namespace to list pods"),
		),
	)

	// Add tool handler
	s.AddTool(tool, listPods)
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func listPods(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["namespace"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}
	gvk := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	resource, err := kubeConfig.getResource(&gvk)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	dynamicclient, err := kubeConfig.dynamicClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ret, err := getPods(dynamicclient, resource, name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(ret), nil
}

func getPods(client *dynamic.DynamicClient, resource *schema.GroupVersionResource, name string) (string, error) {
	podlists, err := client.Resource(*resource).Namespace(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	var podnames []string
	for _, item := range podlists.Items {
		podnames = append(podnames, item.GetName())
	}

	return strings.Join(podnames, ","), nil
}

