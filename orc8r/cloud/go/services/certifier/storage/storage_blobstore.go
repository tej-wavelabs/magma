/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
)

const (
	// CertifierTableBlobstore is the service-wide blobstore table for certifier data
	CertifierTableBlobstore = "certificate_info_blobstore"

	// CertInfoType is the type of CertInfo used in blobstore type fields.
	CertInfoType = "certificate_info"

	// Blobstore needs a network ID, but certifier is network-agnostic so we
	// will use a placeholder value.
	placeholderNetworkID = "placeholder_network"
)

type certifierBlobstore struct {
	factory blobstore.BlobStorageFactory
}

// NewCertifierBlobstore returns an initialized instance of certifierBlobstore as CertifierStorage.
func NewCertifierBlobstore(factory blobstore.BlobStorageFactory) CertifierStorage {
	return &certifierBlobstore{factory: factory}
}

func (c *certifierBlobstore) ListSerialNumbers() ([]string, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	serialNumbers, err := blobstore.ListKeys(store, placeholderNetworkID, CertInfoType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list keys")
	}

	return serialNumbers, store.Commit()
}

func (c *certifierBlobstore) GetCertInfo(serialNumber string) (*protos.CertificateInfo, error) {
	infos, err := c.GetManyCertInfo([]string{serialNumber})
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		return info, nil
	}
	return nil, merrors.ErrNotFound
}

func (c *certifierBlobstore) GetManyCertInfo(serialNumbers []string) (map[string]*protos.CertificateInfo, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	tks := storage.MakeTKs(CertInfoType, serialNumbers)
	blobs, err := store.GetMany(placeholderNetworkID, tks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get many certificate info")
	}

	ret := make(map[string]*protos.CertificateInfo)
	for _, blob := range blobs {
		info := &protos.CertificateInfo{}
		err = proto.Unmarshal(blob.Value, info)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cert info")
		}
		ret[blob.Key] = info
	}

	return ret, store.Commit()
}

func (c *certifierBlobstore) GetAllCertInfo() (map[string]*protos.CertificateInfo, error) {
	infos := map[string]*protos.CertificateInfo{}

	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	serialNumbers, err := blobstore.ListKeys(store, placeholderNetworkID, CertInfoType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list keys")
	}

	if len(serialNumbers) == 0 {
		return infos, store.Commit()
	}

	tks := storage.MakeTKs(CertInfoType, serialNumbers)
	blobs, err := store.GetMany(placeholderNetworkID, tks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get many certificate info")
	}

	for _, blob := range blobs {
		info := &protos.CertificateInfo{}
		err = proto.Unmarshal(blob.Value, info)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cert info")
		}
		infos[blob.Key] = info
	}

	return infos, store.Commit()
}

func (c *certifierBlobstore) PutCertInfo(serialNumber string, certInfo *protos.CertificateInfo) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	marshaledCertInfo, err := proto.Marshal(certInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cert info")
	}

	blob := blobstore.Blob{Type: CertInfoType, Key: serialNumber, Value: marshaledCertInfo}
	err = store.CreateOrUpdate(placeholderNetworkID, blobstore.Blobs{blob})
	if err != nil {
		return errors.Wrap(err, "failed to put certificate info")
	}

	return store.Commit()
}

func (c *certifierBlobstore) DeleteCertInfo(serialNumber string) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	tk := storage.TypeAndKey{Type: CertInfoType, Key: serialNumber}
	err = store.Delete(placeholderNetworkID, []storage.TypeAndKey{tk})
	if err != nil {
		return errors.Wrap(err, "failed to delete certificate info")
	}

	return store.Commit()
}
