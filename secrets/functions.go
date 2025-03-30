package secrets

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/compute/metadata"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/slices"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
)

const (
	pathToSecret                    = "projects/%d/secrets/%s"
	pathToLatestVersion             = "projects/%d/secrets/%s/versions/latest"
	pathToNumericVersion            = "projects/%d/secrets/%s/versions/%d"
	validPathPattern                = "^projects/(\\d+)/secrets/([\\w-]+)$"
	validSecretNamePattern          = "^([\\w-]+)$"
	validSecretAnnotationKeyPattern = "^[A-Za-z0-9]{1}[.\\w-]{0,61}[A-Za-z0-9]{1}$"
	validPathWithVersionPattern     = "^projects/(?P<projectid>\\d+)/secrets/(?P<name>[\\w-]+)(?:/versions/(?P<version>\\d+|latest))?$"
)

var projectNumber = must.Must(strconv.Atoi(must.Must(metadata.NumericProjectIDWithContext(context.Background()))))
var validSecretPathRegex *regexp.Regexp = nil
var validSecretNameRegex *regexp.Regexp = nil
var validSecretPathWithVersionRegex *regexp.Regexp = nil
var validSecretAnnotationKeyRegex *regexp.Regexp = nil

var ValueOutOfRange = errorx.RegisterTrait("value out of range")
var secretVersionNotInValidRange = errorx.NewType(errorx.NewNamespace("SECRET_MANAGER"), "SECRET VERSION NOT IN VALID RANGE", ValueOutOfRange)
var secretVersionNotFound = errorx.NewType(errorx.NewNamespace("SECRET_MANAGER"), "SECRET VERSION NOT FOUND", errorx.NotFound())

func init() {
	var err error
	validSecretPathRegex, err = regexp.Compile(validPathPattern)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not compile regular expression %s because %s", validPathPattern, err.Error())
	}
	validSecretNameRegex, err = regexp.Compile(validSecretNamePattern)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not compile regular expression %s because %s", validSecretNamePattern, err.Error())
	}
	validSecretPathWithVersionRegex, err = regexp.Compile(validPathWithVersionPattern)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not compile regular expression %s because %s", validPathWithVersionPattern, err.Error())
	}
	validSecretAnnotationKeyRegex, err = regexp.Compile(validSecretAnnotationKeyPattern)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not compile regular expression %s because %s", validSecretAnnotationKeyPattern, err.Error())
	}
}

const (
	withoutVersion int = -1
	latestVersion  int = 0
)

type NewPathVersionOption func(path *Path)

func WithVersion(version int) NewPathVersionOption {
	if version <= latestVersion {
		log.Fatal().Msgf("Version %d can not be <= 0, use WithoutVersion() or WithLatestVersion() for those cases", version)
	}
	return func(p *Path) {
		p.Version = version
	}
}

func WithoutVersion() NewPathVersionOption {
	return func(p *Path) {
		p.Version = -1
	}
}

func WithLatestVersion() NewPathVersionOption {
	return func(p *Path) {
		p.Version = 0
	}
}

func NewPath(name string, version NewPathVersionOption) Path {
	p := Path{
		ProjectNumber: projectNumber,
		Name:          name,
		Version:       0,
	}
	version(&p)
	return p
}

type Iterable[I secretmanagerpb.SecretVersion | secretmanagerpb.Secret] interface {
	Next() (*I, error)
}

type SecretVersionIterable struct{}

func (svi SecretVersionIterable) Next() (*secretmanagerpb.SecretVersion, error) {
	return svi.Next()
}

type SecretIterable struct{}

func (si SecretIterable) Next() (*secretmanagerpb.Secret, error) {
	return si.Next()
}

func toSeq2[I secretmanagerpb.SecretVersion | secretmanagerpb.Secret](it Iterable[I]) iter.Seq2[*I, error] {
	return func(yield func(*I, error) bool) {
		for {
			resp, err := it.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				if !yield(nil, errs.IterationError.NewWithNoMessage()) {
					return
				}
			}
			if !yield(resp, nil) {
				return
			}
		}
	}
}

func parsePathFrom(sv *secretmanagerpb.SecretVersion) Path {
	matches := validSecretPathWithVersionRegex.FindStringSubmatch(sv.GetName())
	p := Path{
		ProjectNumber: must.ParseInt(matches[must.Must(slices.FindInSlice(validSecretPathWithVersionRegex.SubexpNames(), "projectid"))]),
		Name:          matches[must.Must(slices.FindInSlice(validSecretPathWithVersionRegex.SubexpNames(), "name"))],
		Version:       must.ParseInt(matches[must.Must(slices.FindInSlice(validSecretPathWithVersionRegex.SubexpNames(), "version"))]),
	}
	return p
}

func buildPathToSecretWithVersion(name string, version int) string {
	return fmt.Sprintf(pathToNumericVersion, projectNumber, name, version)
}

func buildPathToSecretWithLatest(name string) string {
	return fmt.Sprintf(pathToLatestVersion, projectNumber, name)
}

func buildPathToSecretWithoutVersion(name string) string {
	return fmt.Sprintf(pathToSecret, projectNumber, name)
}

func GetSecretValueAsString(ctx context.Context, name string) string {
	value := MustGetSecretValue(ctx, name, GetSecretValue)
	return string(value)
}

func getSecret(ctx context.Context, name string) (*secretmanagerpb.Secret, error) {
	path := buildPathToSecretWithLatest(name)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return nil, errorx.IllegalState.New("%s does not match the validPathWithVersionPattern %s", path, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err = client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.GetSecretRequest{
		Name: path,
	}

	return client.GetSecret(ctx, req)
}

func getSecretLatestVersion(ctx context.Context, name string) (int, error) {
	path := buildPathToSecretWithLatest(name)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return 0, fmt.Errorf("%s does not match the validPathWithVersionPattern %s", path, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return 0, errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.GetSecretVersionRequest{
		Name: path,
	}

	result, err := client.GetSecretVersion(ctx, req)
	if err != nil {
		return 0, errorx.DataUnavailable.Wrap(err, "failed to access secret version: %s", req.GetName())
	}

	return parsePathFrom(result).Version, nil
}

// GetSecretValue accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
func GetSecretValue(ctx context.Context, name string) ([]byte, error) {
	path := buildPathToSecretWithLatest(name)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return nil, errorx.IllegalState.New("%s does not match the validPathWithVersionPattern %s", path, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: path,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, errorx.DataUnavailable.Wrap(err, "failed to access secret version: %s", req.GetName())
	}

	return result.Payload.Data, nil
}

func MustGetSecretValue(ctx context.Context, name string, f func(ctx context.Context, name string) ([]byte, error)) []byte {
	if value, err := f(ctx, name); err != nil {
		err = errors.Join(errs.NotFoundError.New("Error trying to retrieve secret: %s", name), err)
		log.Error().Err(err).Msgf(err.Error())
		panic(errorx.EnsureStackTrace(err))
	} else {
		return value
	}
}

func CreateSecret(ctx context.Context, id string) (*secretmanagerpb.Secret, error) {
	if !validSecretNameRegex.MatchString(id) {
		return nil, fmt.Errorf("%s does not match the validSecretNamePattern pattern %s", id, validSecretNamePattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)
	log.Info().Msgf("attempting to create secret %s", buildPathToSecretWithoutVersion(id))
	path := fmt.Sprintf("projects/%d", projectNumber)
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   path,
		SecretId: id,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := client.CreateSecret(ctx, req)
	if err != nil {
		log.Error().Err(err).Msgf("could not create secret: %s", err.Error())
		if strings.Contains(err.Error(), "AlreadyExists") {
			log.Warn().Err(err).Msg(err.Error())
			pathWithoutVersion := buildPathToSecretWithoutVersion(id)
			return getSecret(ctx, pathWithoutVersion)
		}
		return nil, fmt.Errorf("failed to create secret: %v", err)
	}
	log.Info().Msgf("created secret at %s", secret.Name)
	return secret, nil
}

// AddSecretVersion adds a new secret version to the given secret with the provided payload.
func AddSecretVersion(ctx context.Context, name string, value []byte) (*secretmanagerpb.SecretVersion, error) {
	// parent := "projects/my-project/secrets/my-secret"
	path := buildPathToSecretWithoutVersion(name)
	if !validSecretPathRegex.MatchString(path) {
		return nil, fmt.Errorf("%s does not match the validPathPattern %s", path, validPathPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	secretVersionRequest := &secretmanagerpb.AddSecretVersionRequest{
		Parent: path,
		Payload: &secretmanagerpb.SecretPayload{
			Data: value,
		},
	}

	secretVersion, err := client.AddSecretVersion(ctx, secretVersionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to add secret version: %v", err)
	}
	log.Info().Msgf("added version to secret %s", secretVersion.Name)
	return secretVersion, nil
}

func EnableSecretVersion(ctx context.Context, name string, version int) error {
	if version <= 0 {
		return secretVersionNotFound.New("version %d out of range, must be >= 1", version)
	}
	// path := "projects/my-project/secrets/my-secret/versions/5"
	path := buildPathToSecretWithVersion(name, version)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return fmt.Errorf("%s does not match the required pattern %s", name, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.EnableSecretVersionRequest{
		Name: path,
	}

	result, err := client.EnableSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to enable secret version: %v", err)
	}
	log.Info().Msgf("enabled secret at %s", result.Name)
	return nil
}

func DisableSecretVersion(ctx context.Context, name string, version int) error {
	if version <= 0 {
		return secretVersionNotFound.New("version %d out of range, must be >= 1", version)
	}
	// path := "projects/my-project/secrets/my-secret/versions/5"
	path := buildPathToSecretWithVersion(name, version)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return fmt.Errorf("%s does not match the required pattern %s", name, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.DisableSecretVersionRequest{
		Name: path,
	}

	result, err := client.DisableSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to enable secret version: %v", err)
	}
	log.Info().Msgf("disabled secret at %s", result.Name)
	return nil
}

func DestroySecretVersion(ctx context.Context, name string, version int) error {
	if version <= 0 {
		return secretVersionNotFound.New("version %d out of range, must be >= 1", version)
	}
	// path := "projects/my-project/secrets/my-secret/versions/5"
	path := buildPathToSecretWithVersion(name, version)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return fmt.Errorf("%s does not match the required pattern %s", name, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.DestroySecretVersionRequest{
		Name: path,
	}

	_, err = client.DestroySecretVersion(ctx, req)
	if err != nil {
		if err.Error() == "not found" {
			log.Warn().Err(err).Msgf("failed to destroy secret version: %v", err)
			return nil
		}
		return fmt.Errorf("failed to destroy secret version: %v", err)
	}
	return nil
}

func DestroyAllButLatestVersion(ctx context.Context, name string) error {
	version, err := getSecretLatestVersion(ctx, name)
	if err != nil {
		return err
	}
	return DestroyAllPreviousVersions(ctx, name, version)
}

func DestroyAllPreviousVersions(ctx context.Context, name string, version int) error {
	// path := "projects/${project-number}/secrets/${name}/versions/${version}"
	path := buildPathToSecretWithoutVersion(name)
	if !validSecretPathWithVersionRegex.MatchString(path) {
		return fmt.Errorf("%s does not match the required pattern %s", name, validPathWithVersionPattern)
	}
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err = client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)

	req := &secretmanagerpb.ListSecretVersionsRequest{
		Parent: path,
	}
	if version == 0 {
		version, err = getSecretLatestVersion(ctx, name)
		if err != nil {
			return err
		}
	}
	sviter := toSeq2[secretmanagerpb.SecretVersion](client.ListSecretVersions(ctx, req))
	//log.Debug().Msgf("Secret Version Iterator.PageInfo(): %s", sviter.PageInfo())
	for sv, iterr := range sviter {
		if iterr != nil {
			log.Error().Err(iterr).Msg(iterr.Error())
			err = errors.Join(err, iterr)
		} else {
			p := parsePathFrom(sv)
			if p.Version < version && version > 0 {
				err = DestroySecretVersion(ctx, name, p.Version)
				if err != nil {
					return errs.NotDeletedError.WrapWithNoMessage(iterr)
				}
			}
		}
	}
	return err
}

func UpdateSecretWithNewVersion(ctx context.Context, name string, value []byte) error {
	newSecretVersion, err := AddSecretVersion(ctx, name, value)
	if err != nil {
		return err
	}
	log.Info().Msgf("created new secret version %s", newSecretVersion.Name)
	path := parsePathFrom(newSecretVersion)
	return EnableSecretVersion(ctx, name, path.Version)
}

func ReplaceSecretWithNewVersion(ctx context.Context, name string, value []byte) error {
	newSecretVersion, err := AddSecretVersion(ctx, name, value)
	if err != nil {
		return err
	}
	log.Info().Msgf("created new secret version %s", newSecretVersion.Name)
	path := parsePathFrom(newSecretVersion)
	if path.Version > 1 {
		return DestroySecretVersion(ctx, name, path.Version-1)
	}
	return nil
}

func CreateSecretWithValue(ctx context.Context, name string, value []byte) error {
	log.Info().Msgf("creating secret with value and enabling at %s", buildPathToSecretWithoutVersion(name))
	_, err := CreateSecret(ctx, name)
	if err != nil {
		return err
	}
	return UpdateSecretWithNewVersion(ctx, name, value)
}

func RemoveSecret(ctx context.Context, name string) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "failed to create secretmanager client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
		}
	}(client)
	path := buildPathToSecretWithoutVersion(name)
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: path,
	}
	return client.DeleteSecret(ctx, req)
}
