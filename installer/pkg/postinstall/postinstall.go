package postinstall

func DownloadPostInstallPackages() error {
	packages, err := getPackageList()
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := downloadAllPackages(packages); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}
