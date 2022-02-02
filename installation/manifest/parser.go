package manifest

import (
	"encoding/pem"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	biproperty "github.com/cloudfoundry/bosh-utils/property"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/cppforlife/go-patch/patch"
	"gopkg.in/yaml.v2"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	birelsetmanifest "github.com/cloudfoundry/bosh-cli/release/set/manifest"
)

type Parser interface {
	Parse(string, boshtpl.Variables, patch.Op, birelsetmanifest.Manifest) (Manifest, error)
}

type parser struct {
	fs            boshsys.FileSystem
	uuidGenerator boshuuid.Generator
	logger        boshlog.Logger
	logTag        string
	validator     Validator
}

type manifest struct {
	Name          string
	CloudProvider installation `yaml:"cloud_provider"`
}

type installation struct {
	Template   template
	Properties map[interface{}]interface{}
	Mbus       string
	Cert       Certificate
}

type template struct {
	Name    string
	Release string
}

func NewParser(fs boshsys.FileSystem, uuidGenerator boshuuid.Generator, logger boshlog.Logger, validator Validator) Parser {
	return &parser{
		fs:            fs,
		uuidGenerator: uuidGenerator,
		logger:        logger,
		logTag:        "deploymentParser",
		validator:     validator,
	}
}

func (p *parser) Parse(path string, vars boshtpl.Variables, op patch.Op, releaseSetManifest birelsetmanifest.Manifest) (Manifest, error) {
	contents, err := p.fs.ReadFile(path)
	if err != nil {
		return Manifest{}, bosherr.WrapErrorf(err, "Reading file %s", path)
	}

	tpl := boshtpl.NewTemplate(contents)

	bytes, err := tpl.Evaluate(vars, op, boshtpl.EvaluateOpts{ExpectAllKeys: true})
	if err != nil {
		return Manifest{}, bosherr.WrapErrorf(err, "Evaluating manifest")
	}

	comboManifest := manifest{}

	err = yaml.Unmarshal(bytes, &comboManifest)
	if err != nil {
		return Manifest{}, bosherr.WrapError(err, "Unmarshalling installation manifest")
	}

	p.logger.Debug(p.logTag, "Parsed installation manifest: %#v", comboManifest)

	if comboManifest.CloudProvider.Cert.CA != "" {
		pkey, _ := pem.Decode([]byte(comboManifest.CloudProvider.Cert.CA))
		if pkey == nil {
			return Manifest{}, bosherr.Error("Invalid CA cert")
		}
	}

	installationManifest := Manifest{
		Name: comboManifest.Name,
		Template: ReleaseJobRef{
			Name:    comboManifest.CloudProvider.Template.Name,
			Release: comboManifest.CloudProvider.Template.Release,
		},
		Mbus: comboManifest.CloudProvider.Mbus,
		Cert: comboManifest.CloudProvider.Cert,
	}

	properties, err := biproperty.BuildMap(comboManifest.CloudProvider.Properties)
	if err != nil {
		return Manifest{}, bosherr.WrapErrorf(err, "Parsing cloud_provider manifest properties: %#v", comboManifest.CloudProvider.Properties)
	}
	installationManifest.Properties = properties

	err = p.validator.Validate(installationManifest, releaseSetManifest)
	if err != nil {
		return Manifest{}, bosherr.WrapError(err, "Validating installation manifest")
	}

	return installationManifest, nil
}

func (p *parser) lookForPrivateSshHeader(key string) bool {
	return strings.HasPrefix(key, "-----BEGIN RSA PRIVATE KEY-----") ||
		strings.HasPrefix(key, "-----BEGIN OPENSSH PRIVATE KEY-----")
}
