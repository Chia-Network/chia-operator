/*
Copyright 2025 Chia Network Inc.
*/

package chiaca

import (
	"context"
	"fmt"
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MockClient is a mock implementation of the client.Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	args := m.Called(ctx, key, obj, opts)
	return args.Error(0)
}

func (m *MockClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	args := m.Called(ctx, list, opts)
	return args.Error(0)
}

func (m *MockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *MockClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *MockClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *MockClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	args := m.Called(ctx, obj, patch, opts)
	return args.Error(0)
}

func (m *MockClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *MockClient) Status() client.StatusWriter {
	args := m.Called()
	return args.Get(0).(client.StatusWriter)
}

func (m *MockClient) FieldIndexer() client.FieldIndexer {
	args := m.Called()
	return args.Get(0).(client.FieldIndexer)
}

func (m *MockClient) SubResource(subResource string) client.SubResourceClient {
	args := m.Called(subResource)
	return args.Get(0).(client.SubResourceClient)
}

func (m *MockClient) Scheme() *runtime.Scheme {
	args := m.Called()
	return args.Get(0).(*runtime.Scheme)
}

func (m *MockClient) RESTMapper() meta.RESTMapper {
	args := m.Called()
	return args.Get(0).(meta.RESTMapper)
}

func (m *MockClient) WithWatch() client.WithWatch {
	args := m.Called()
	return args.Get(0).(client.WithWatch)
}

func (m *MockClient) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	args := m.Called(obj)
	return args.Get(0).(schema.GroupVersionKind), args.Error(1)
}

func (m *MockClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	args := m.Called(obj)
	return args.Bool(0), args.Error(1)
}

// MockStatusWriter is a mock implementation of the client.StatusWriter interface
type MockStatusWriter struct {
	mock.Mock
}

func (m *MockStatusWriter) Create(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	args := m.Called(ctx, obj, subResource, opts)
	return args.Error(0)
}

func (m *MockStatusWriter) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *MockStatusWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	args := m.Called(ctx, obj, patch, opts)
	return args.Error(0)
}

var testChiaCA = k8schianetv1.ChiaCA{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ChiaCA",
		APIVersion: "k8s.chia.net/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testname",
		Namespace: "testnamespace",
	},
}

func TestGetChiaCASecretName_DefaultName(t *testing.T) {
	// Test with default name (no custom secret name specified)
	secretName := getChiaCASecretName(testChiaCA)
	assert.Equal(t, "testname", secretName)
}

func TestGetChiaCASecretName_CustomName(t *testing.T) {
	// Test with custom secret name
	customCA := testChiaCA
	customCA.Spec.Secret = "custom-secret-name"
	secretName := getChiaCASecretName(customCA)
	assert.Equal(t, "custom-secret-name", secretName)
}

func TestGetChiaCASecretName_EmptyString(t *testing.T) {
	// Test with empty string (should use default name)
	customCA := testChiaCA
	customCA.Spec.Secret = ""
	secretName := getChiaCASecretName(customCA)
	assert.Equal(t, "testname", secretName)
}

func TestGetChiaCASecretName_WhitespaceString(t *testing.T) {
	// Test with whitespace string (should use default name)
	customCA := testChiaCA
	customCA.Spec.Secret = "   "
	secretName := getChiaCASecretName(customCA)
	assert.Equal(t, "testname", secretName)
}

func TestCASecretExists_SecretExists(t *testing.T) {
	// Create a mock client
	mockClient := new(MockClient)

	// Set up the mock to return no error (secret exists)
	mockClient.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create a reconciler with the mock client
	reconciler := &ChiaCAReconciler{
		Client: mockClient,
	}

	// Call the function
	exists, err := reconciler.caSecretExists(context.Background(), testChiaCA)

	// Assert the results
	assert.NoError(t, err)
	assert.True(t, exists)
	mockClient.AssertExpectations(t)
}

func TestCASecretExists_SecretNotFound(t *testing.T) {
	// Create a mock client
	mockClient := new(MockClient)

	// Set up the mock to return a NotFound error
	mockClient.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.NewNotFound(schema.GroupResource{Group: "", Resource: "secrets"}, "testname"))

	// Create a reconciler with the mock client
	reconciler := &ChiaCAReconciler{
		Client: mockClient,
	}

	// Call the function
	exists, err := reconciler.caSecretExists(context.Background(), testChiaCA)

	// Assert the results
	assert.NoError(t, err)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}

func TestCASecretExists_OtherError(t *testing.T) {
	// Create a mock client
	mockClient := new(MockClient)

	// Set up the mock to return a different error
	mockClient.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("some other error"))

	// Create a reconciler with the mock client
	reconciler := &ChiaCAReconciler{
		Client: mockClient,
	}

	// Call the function
	exists, err := reconciler.caSecretExists(context.Background(), testChiaCA)

	// Assert the results
	assert.Error(t, err)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}
