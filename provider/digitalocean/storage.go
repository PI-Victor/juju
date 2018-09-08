package digitalocean

import (
	"github.com/juju/errors"

	envCtx "github.com/juju/juju/environs/context"
	"github.com/juju/juju/storage"
)

type doStorageProvider struct{}

func (s *doStorageProvider) StorageProvider(t storage.ProviderType) (storage.Provider, error) {
	return &doStorageProvider{}, nil
}

// Supports is part of the Provider interface.
func (s *doStorageProvider) Supports(k storage.StorageKind) bool {
	return false
}

func (e *doStorageProvider) ValidateConfig(cfg *storage.Config) error {
	return nil
}

// Scope is part of the Provider interface.
func (s *doStorageProvider) Scope() storage.Scope {
	return storage.ScopeEnviron
}

// Dynamic is part of the Provider interface.
func (e *doStorageProvider) Dynamic() bool {
	return false
}

// Releasable is part of the Provider interface.
func (e *doStorageProvider) Releasable() bool {
	// NOTE(axw) Azure storage is currently tied to a model, and cannot
	// be released or imported. To support releasing and importing, we'll
	// need Azure to support moving managed disks between resource groups.
	return false
}

// DefaultPools is part of the Provider interface.
func (e *doStorageProvider) DefaultPools() []*storage.Config {
	return nil
}

// VolumeSource is part of the Provider interface.
func (e *doStorageProvider) VolumeSource(cfg *storage.Config) (storage.VolumeSource, error) {
	return &volumeSource{}, nil
}

// FilesystemSource is part of the Provider interface.
func (e *doStorageProvider) FilesystemSource(providerConfig *storage.Config) (storage.FilesystemSource, error) {
	return nil, errors.NotSupportedf("filesystems")
}

type volumeSource struct{}

func (v *volumeSource) CreateVolumes(ctx envCtx.ProviderCallContext, params []storage.VolumeParams) ([]storage.CreateVolumesResult, error) {
	return nil, nil
}

func (v *volumeSource) ListVolumes(ctx envCtx.ProviderCallContext) ([]string, error) {
	return nil, nil
}

// DescribeVolumes returns the properties of the volumes with the
// specified provider volume IDs.
func (v *volumeSource) DescribeVolumes(ctx envCtx.ProviderCallContext, volIds []string) ([]storage.DescribeVolumesResult, error) {
	return nil, nil
}

// DestroyVolumes destroys the volumes with the specified provider
// volume IDs.
func (v *volumeSource) DestroyVolumes(ctx envCtx.ProviderCallContext, volIds []string) ([]error, error) {
	return nil, nil
}

// ReleaseVolumes releases the volumes with the specified provider
// volume IDs from the model/controller.
func (v *volumeSource) ReleaseVolumes(ctx envCtx.ProviderCallContext, volIds []string) ([]error, error) {
	return nil, nil
}

// ValidateVolumeParams validates the provided volume creation
// parameters, returning an error if they are invalid.
func (v *volumeSource) ValidateVolumeParams(params storage.VolumeParams) error {
	return nil
}

func (v *volumeSource) AttachVolumes(ctx envCtx.ProviderCallContext, params []storage.VolumeAttachmentParams) ([]storage.AttachVolumesResult, error) {
	return nil, nil
}

func (v *volumeSource) DetachVolumes(ctx envCtx.ProviderCallContext, params []storage.VolumeAttachmentParams) ([]error, error) {
	return nil, nil
}
