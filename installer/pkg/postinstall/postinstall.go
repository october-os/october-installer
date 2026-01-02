package postinstall

func InstallPostInstallPackages() error {
	packages, err := getPackageList(packageFilePath)
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return downloadAllPackages(packages, false)
}

func InstallAurHelperAndPackages() error {
	if err := activateBuilderAccount(); err != nil {
		return err
	}

	if err := installYay(); err != nil {
		return err
	}

	packages, err := getPackageList(aurFilePath)
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := downloadAllPackages(packages, true); err != nil {
		return err
	}

	return deleteBuilderAccount()
}
