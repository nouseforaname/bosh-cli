package manifest_test

import (
	"errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	biproperty "github.com/cloudfoundry/bosh-utils/property"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	fakeuuid "github.com/cloudfoundry/bosh-utils/uuid/fakes"
	"github.com/cppforlife/go-patch/patch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cloudfoundry/bosh-cli/installation/manifest"
	"github.com/cloudfoundry/bosh-cli/installation/manifest/fakes"
	birelsetmanifest "github.com/cloudfoundry/bosh-cli/release/set/manifest"
)

type manifestFixtures struct {
	validManifest             string
	missingPrivateKeyManifest string
}

var _ = Describe("Parser", func() {
	comboManifestPath := "/path/to/fake-deployment-manifest"
	releaseSetManifest := birelsetmanifest.Manifest{}
	var (
		fakeFs            *fakesys.FakeFileSystem
		fakeUUIDGenerator *fakeuuid.FakeGenerator
		parser            manifest.Parser
		logger            boshlog.Logger
		fakeValidator     *fakes.FakeValidator
		fixtures          manifestFixtures
	)
	BeforeEach(func() {
		fakeValidator = fakes.NewFakeValidator()
		fakeValidator.SetValidateBehavior([]fakes.ValidateOutput{
			{Err: nil},
		})
		fakeFs = fakesys.NewFakeFileSystem()
		logger = boshlog.NewLogger(boshlog.LevelNone)
		fakeUUIDGenerator = fakeuuid.NewFakeGenerator()
		parser = manifest.NewParser(fakeFs, fakeUUIDGenerator, logger, fakeValidator)
		fixtures = manifestFixtures{
			validManifest: `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
  properties:
    fake-property-name:
      nested-property: fake-property-value
`,
			missingPrivateKeyManifest: `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`,
		}
	})

	Describe("#Parse", func() {
		Context("when combo manifest path does not exist", func() {
			It("returns an error", func() {
				_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when parser fails to read the combo manifest file", func() {
			JustBeforeEach(func() {
				fakeFs.WriteFileString(comboManifestPath, fixtures.validManifest)
				fakeFs.ReadFileError = errors.New("fake-read-file-error")
			})

			It("returns an error", func() {
				_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with a valid manifest", func() {
			BeforeEach(func() {
				fakeFs.WriteFileString(comboManifestPath, fixtures.validManifest)
			})

			It("parses installation from combo manifest", func() {
				installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).ToNot(HaveOccurred())

				Expect(installationManifest).To(Equal(manifest.Manifest{
					Name: "fake-deployment-name",
					Template: manifest.ReleaseJobRef{
						Name:    "fake-cpi-job-name",
						Release: "fake-cpi-release-name",
					},
					Properties: biproperty.Map{
						"fake-property-name": biproperty.Map{
							"nested-property": "fake-property-value",
						},
					},
					Mbus: "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
				}))
			})
		})

		Context("when ssh tunnel config is present", func() {
			Context("with raw private key", func() {
				Context("that is valid", func() {
					BeforeEach(func() {
						fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("sets the raw private key field", func() {
						installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).ToNot(HaveOccurred())

						Expect(installationManifest).To(Equal(manifest.Manifest{
							Name: "fake-deployment-name",
							Template: manifest.ReleaseJobRef{
								Name:    "fake-cpi-job-name",
								Release: "fake-cpi-release-name",
							},
							Properties: biproperty.Map{},
							Mbus:       "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
						}))

					})
				})
				Context("that is invalid", func() {
					BeforeEach(func() {
						fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})
				})
			})

			Context("with new format raw private key", func() {
				Context("that is valid", func() {
					BeforeEach(func() {
						fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("sets the raw private key field", func() {
						installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).ToNot(HaveOccurred())

						Expect(installationManifest).To(Equal(manifest.Manifest{
							Name: "fake-deployment-name",
							Template: manifest.ReleaseJobRef{
								Name:    "fake-cpi-job-name",
								Release: "fake-cpi-release-name",
							},
							Properties: biproperty.Map{},
							Mbus:       "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
						}))

					})
				})
				Context("that is invalid", func() {
					BeforeEach(func() {
						fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

				})
			})

		})

		It("handles installation manifest validation errors", func() {
			fakeFs.WriteFileString(comboManifestPath, fixtures.validManifest)

			fakeValidator.SetValidateBehavior([]fakes.ValidateOutput{
				{Err: errors.New("nope")},
			})

			_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Validating installation manifest: nope"))
		})

		Context("when interpolating variables", func() {
			BeforeEach(func() {
				fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
				fakeFs.ExpandPathExpanded = "/Users/foo/tmp/fake-ssh-key.pem"

				fakeFs.WriteFileString(comboManifestPath, `---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: ((url))
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
				fakeFs.WriteFileString("/Users/foo/tmp/fake-ssh-key.pem", "--- BEGIN KEY --- blah --- END KEY ---")
			})

			It("resolves their values", func() {
				vars := boshtpl.StaticVariables{"url": "~/tmp/fake-ssh-key.pem"}
				ops := patch.Ops{
					patch.ReplaceOp{Path: patch.MustNewPointerFromString("/name"), Value: "replaced-name"},
				}

				installationManifest, err := parser.Parse(comboManifestPath, vars, ops, releaseSetManifest)
				Expect(err).ToNot(HaveOccurred())

				Expect(installationManifest).To(Equal(manifest.Manifest{
					Name: "replaced-name",
					Template: manifest.ReleaseJobRef{
						Name:    "fake-cpi-job-name",
						Release: "fake-cpi-release-name",
					},
					Properties: biproperty.Map{},
					Mbus:       "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
				}))
			})

			It("returns an error if variable key is missing", func() {
				_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Expected to find variables: url"))
			})
		})

		Context("when CA cert is present", func() {
			Context("with raw certificate", func() {
				Context("that is valid", func() {
					BeforeEach(func() {
						fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
  cert:
    ca: |
      -----BEGIN CERTIFICATE-----
      MIIC+TCCAeGgAwIBAgIQLzf5Fs3v+Dblm+CKQFxiKTANBgkqhkiG9w0BAQsFADAm
      MQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoTDUNsb3VkIEZvdW5kcnkwHhcNMTcwNTE2
      MTUzNTI4WhcNMTgwNTE2MTUzNTI4WjAmMQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoT
      DUNsb3VkIEZvdW5kcnkwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC+
      4E0QJMOpQwbHACvrZ4FleP4/DMFvYUBySfKzDOgd99Nm8LdXuJcI1SYHJ3sV+mh0
      +cQmRt8U2A/lw7bNU6JdM0fWHa/2nGjSBKWgPzba68NdsmwjqUjLatKpr1yvd384
      PJJKC7NrxwvChgB8ui84T4SrXHCioYMDEDIqLGmHJHMKnzQ17nu7ECO4e6QuCfnH
      RDs7dTjomTAiFuF4fh4SPgEDMGaCE5HZr4t3gvc9n4UftpcCpi+Jh+neRiWx+v37
      ZAYf2kp3wWtYDlgWk06cZzHZZ9uYZFwHDNHdDKHxGGvAh2Rm6rpPF2oA6OEyx6BH
      85/STCgSMCnV1Wkd+1yPAgMBAAGjIzAhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMB
      Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBGvGggx3IM4KCMpVDSv9zFKX4K
      IuCRQ6VFab3sgnlelMFaMj3+8baJ/YMko8PP1wVfUviVgKuiZO8tqL00Yo4s1WKp
      x3MLIG4eBX9pj0ZVRa3kpcF2Wvg6WhrzUzONf7pfuz/9avl77o4aSt4TwyCvM4Iu
      gJ7quVQKcfQcAVwuwWRrZXyhjhHaVKoPP5yRS+ESVTl70J5HBh6B7laooxf1yVAW
      8NJK1iQ1Pw2x3ABBo1cSMcTQ3Hk1ZWThJ7oPul2+QyzvOjIjiEPBstyzEPaxPG4I
      nH9ttalAwSLBsobVaK8mmiAdtAdx+CmHWrB4UNxCPYasrt5A6a9A9SiQ2dLd
      -----END CERTIFICATE-----
`)
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("sets the CA cert field", func() {
						installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).ToNot(HaveOccurred())

						Expect(installationManifest.Cert.CA).To(Equal(`-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIQLzf5Fs3v+Dblm+CKQFxiKTANBgkqhkiG9w0BAQsFADAm
MQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoTDUNsb3VkIEZvdW5kcnkwHhcNMTcwNTE2
MTUzNTI4WhcNMTgwNTE2MTUzNTI4WjAmMQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoT
DUNsb3VkIEZvdW5kcnkwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC+
4E0QJMOpQwbHACvrZ4FleP4/DMFvYUBySfKzDOgd99Nm8LdXuJcI1SYHJ3sV+mh0
+cQmRt8U2A/lw7bNU6JdM0fWHa/2nGjSBKWgPzba68NdsmwjqUjLatKpr1yvd384
PJJKC7NrxwvChgB8ui84T4SrXHCioYMDEDIqLGmHJHMKnzQ17nu7ECO4e6QuCfnH
RDs7dTjomTAiFuF4fh4SPgEDMGaCE5HZr4t3gvc9n4UftpcCpi+Jh+neRiWx+v37
ZAYf2kp3wWtYDlgWk06cZzHZZ9uYZFwHDNHdDKHxGGvAh2Rm6rpPF2oA6OEyx6BH
85/STCgSMCnV1Wkd+1yPAgMBAAGjIzAhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMB
Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBGvGggx3IM4KCMpVDSv9zFKX4K
IuCRQ6VFab3sgnlelMFaMj3+8baJ/YMko8PP1wVfUviVgKuiZO8tqL00Yo4s1WKp
x3MLIG4eBX9pj0ZVRa3kpcF2Wvg6WhrzUzONf7pfuz/9avl77o4aSt4TwyCvM4Iu
gJ7quVQKcfQcAVwuwWRrZXyhjhHaVKoPP5yRS+ESVTl70J5HBh6B7laooxf1yVAW
8NJK1iQ1Pw2x3ABBo1cSMcTQ3Hk1ZWThJ7oPul2+QyzvOjIjiEPBstyzEPaxPG4I
nH9ttalAwSLBsobVaK8mmiAdtAdx+CmHWrB4UNxCPYasrt5A6a9A9SiQ2dLd
-----END CERTIFICATE-----
`))
					})
				})

				Context("that is invalid", func() {
					BeforeEach(func() {
						fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
  cert:
    ca: |
      -----BEGIN CERTIFICATE-----
      no valid certificate
      -----END CERTIFICATE-----
`)
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Invalid CA cert"))
					})
				})
			})

			Context("when ca cert is not provided", func() {
				BeforeEach(func() {
					fakeFs.WriteFileString(comboManifestPath, fixtures.missingPrivateKeyManifest)
				})

				It("does not expand the path", func() {
					installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
					Expect(err).ToNot(HaveOccurred())

					Expect(installationManifest.Cert.CA).To(Equal(""))
				})
			})
		})
	})
})
