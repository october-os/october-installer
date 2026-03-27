package hostname

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValidHostname(t *testing.T) {
	hostname := "validhostname"

	err := ValidateHostname(hostname)

	assert.NoError(t, err, "Hostname should be valid")
}

func TestValidateEmptyHostname(t *testing.T) {
	hostname := ""

	err := ValidateHostname(hostname)

	assert.Error(t, err, "Hostname shouldn't be empty")
	assert.ErrorAs(t, err, &HostnameError{})
}

func TestValidateTooLongHostname(t *testing.T) {
	hostname := " Loremipsumdolorsitametconsecteturadipiscingelitligula"

	err := ValidateHostname(hostname)

	assert.Error(t, err, "Hostname can't be longer than 63 characters")
	assert.ErrorAs(t, err, &HostnameError{})
}

func TestValidateAllCapsHostname(t *testing.T) {
	hostname := "HELLOWORLD"

	err := ValidateHostname(hostname)

	assert.Error(t, err, "Hostname can't contain uppercase letters")
	assert.ErrorAs(t, err, &HostnameError{})
}

func TestValidateHostnameStartingWithDash(t *testing.T) {
	hostname := "-helloworld"

	err := ValidateHostname(hostname)

	assert.Error(t, err, "Hostname can't starts with a dash")
	assert.ErrorAs(t, err, &HostnameError{})
}

func TestValidateHostnameWithSpaces(t *testing.T) {
	hostname := "hello world"

	err := ValidateHostname(hostname)

	assert.Error(t, err, "Hostname can't contain spaces")
	assert.ErrorAs(t, err, &HostnameError{})
}

func TestSetHostname(t *testing.T) {
	hostname := "testhostname"
	originalFilePath := hostnameFilePath

	tempDir := t.TempDir()
	tempFile, _ := os.CreateTemp(tempDir, "hostname-test")
	hostnameFilePath = tempFile.Name()

	defer func() {
		hostnameFilePath = originalFilePath
	}()

	err := SetHostname(hostname)

	got, _ := os.ReadFile(tempFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, hostname, string(got))
}
