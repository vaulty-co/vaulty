package proxy

import "testing"

func TestVaultID(t *testing.T) {
	vaultID, err := getVaultID("vlt123.proxy.test")
	if err != nil {
		t.Error(err)
	}

	if vaultID != "vlt123" {
		t.Errorf("vaultID = %s; want vlt123", vaultID)
	}
}
