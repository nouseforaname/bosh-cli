package instance

import (
	biblobstore "github.com/cloudfoundry/bosh-cli/blobstore"
	bicloud "github.com/cloudfoundry/bosh-cli/cloud"
	bivm "github.com/cloudfoundry/bosh-cli/deployment/vm"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type ManagerFactory interface {
	NewManager(bicloud.Cloud, bivm.Manager, biblobstore.Blobstore) Manager
}

type managerFactory struct {
	instanceFactory Factory
	logger          boshlog.Logger
}

func NewManagerFactory(
	instanceFactory Factory,
	logger boshlog.Logger,
) ManagerFactory {
	return &managerFactory{
		instanceFactory: instanceFactory,
		logger:          logger,
	}
}

func (f *managerFactory) NewManager(cloud bicloud.Cloud, vmManager bivm.Manager, blobstore biblobstore.Blobstore) Manager {
	return NewManager(
		cloud,
		vmManager,
		blobstore,
		f.instanceFactory,
		f.logger,
	)
}
