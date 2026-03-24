package core

import (
	"testing"

	"github.com/klauspost/cpuid/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetIntelCpuMicrocode(t *testing.T) {
	want := intelMicrocode
	cpuid.CPU.VendorID = cpuid.Intel

	got, err := getCpuMicroCode()

	assert.Nil(t, err, "Error should be nil with valid CPU")
	assert.Equal(t, want, got, "Should be intel microcode if detects intel CPU")
}

func TestGetAmdCpuMicrocode(t *testing.T) {
	want := amdMicrocode
	cpuid.CPU.VendorID = cpuid.AMD

	got, err := getCpuMicroCode()

	assert.Nil(t, err, "Error should be nil with valid CPU")
	assert.Equal(t, want, got, "Should be amd microcode if detects amd CPU")
}

func TestGetInvalidCpuMicrocode(t *testing.T) {
	cpuid.CPU.VendorID = cpuid.Apple

	got, err := getCpuMicroCode()

	assert.Equal(t, "", got, "Invalid CPU should return an empty string")
	assert.NotNil(t, err, "Error should not be nil on invalid CPU")
}
