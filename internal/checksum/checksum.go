package checksum

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

const (
	AlgorithmSHA256 = "sha256"
	AlgorithmSHA512 = "sha512"
)

type Spec struct {
	Algorithm string
	Expected  string
}

func Parse(raw string) (Spec, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Spec{}, fmt.Errorf("checksum cannot be empty")
	}

	algorithm, expected, ok := strings.Cut(trimmed, ":")
	if !ok {
		return Spec{}, fmt.Errorf("checksum must use format <algorithm>:<hex>")
	}

	algorithm = strings.ToLower(strings.TrimSpace(algorithm))
	expected = strings.ToLower(strings.TrimSpace(expected))
	if algorithm == "" || expected == "" {
		return Spec{}, fmt.Errorf("checksum must use format <algorithm>:<hex>")
	}

	length, err := expectedLength(algorithm)
	if err != nil {
		return Spec{}, err
	}
	if len(expected) != length {
		return Spec{}, fmt.Errorf("%s checksum must be %d hex characters", algorithm, length)
	}
	if _, err := hex.DecodeString(expected); err != nil {
		return Spec{}, fmt.Errorf("%s checksum must be valid hexadecimal", algorithm)
	}

	return Spec{
		Algorithm: algorithm,
		Expected:  expected,
	}, nil
}

func (s Spec) String() string {
	if s.Algorithm == "" || s.Expected == "" {
		return ""
	}
	return s.Algorithm + ":" + s.Expected
}

func ComputeFile(path, algorithm string) (string, error) {
	hasher, err := newHasher(algorithm)
	if err != nil {
		return "", err
	}

	// #nosec G304 -- checksum verification opens the completed downloader target path after target safety validation.
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open checksum target %q: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("read checksum target %q: %w", path, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func VerifyFile(path string, spec Spec) (string, error) {
	actual, err := ComputeFile(path, spec.Algorithm)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(spec.Expected, actual) {
		return actual, fmt.Errorf("checksum mismatch: expected %s, got %s", spec.Expected, actual)
	}
	return actual, nil
}

func SupportedAlgorithms() []string {
	return []string{AlgorithmSHA256, AlgorithmSHA512}
}

func expectedLength(algorithm string) (int, error) {
	switch algorithm {
	case AlgorithmSHA256:
		return sha256.Size * 2, nil
	case AlgorithmSHA512:
		return sha512.Size * 2, nil
	default:
		return 0, fmt.Errorf("unsupported checksum algorithm %q", algorithm)
	}
}

func newHasher(algorithm string) (hash.Hash, error) {
	switch strings.ToLower(strings.TrimSpace(algorithm)) {
	case AlgorithmSHA256:
		return sha256.New(), nil
	case AlgorithmSHA512:
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unsupported checksum algorithm %q", algorithm)
	}
}
