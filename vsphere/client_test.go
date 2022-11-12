package vsphere

import (
	"testing"
)

func TestNewVSphereClient(t *testing.T) {
	NewVSphereClient(false)
}

func TestSetCredential(t *testing.T) {
	vsc := NewVSphereClient(false)
	vsc.SetCredential("fake-user", "fake-password")
}

func TestValidatePath(t *testing.T) {
	test_path := []string{
		"/api/",
		"/api/session",
		"/api/vcenter",
		"/api/vcenter/vm",
	}
	vsc := NewVSphereClient(false)
	for _, path := range test_path {
		err := vsc.validatePath(path)
		if err != nil {
			t.Errorf("path = '%v', want /api/*", path)
		}
	}
}
