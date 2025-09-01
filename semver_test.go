package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSemantic(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
		expected    *SemVer
	}{
		{
			name:        "valid version",
			version:     "1.2.3",
			expectError: false,
			expected: &SemVer{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
		},
		{
			name:        "valid version with v prefix",
			version:     "v1.2.3",
			expectError: false,
			expected: &SemVer{
				major:    1,
				minor:    2,
				patch:    3,
				original: "v1.2.3",
			},
		},
		{
			name:        "version with prerelease",
			version:     "1.2.3-alpha",
			expectError: false,
			expected: &SemVer{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha",
				original:   "1.2.3-alpha",
			},
		},
		{
			name:        "version with prerelease and metadata",
			version:     "1.2.3-alpha.1+build.123",
			expectError: false,
			expected: &SemVer{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha.1",
				metadata:   "build.123",
				original:   "1.2.3-alpha.1+build.123",
			},
		},
		{
			name:        "version with metadata only",
			version:     "1.2.3+build.123",
			expectError: false,
			expected: &SemVer{
				major:    1,
				minor:    2,
				patch:    3,
				metadata: "build.123",
				original: "1.2.3+build.123",
			},
		},
		{
			name:        "empty version",
			version:     "",
			expectError: true,
		},
		{
			name:        "invalid format",
			version:     "1.2",
			expectError: true,
		},
		{
			name:        "invalid major",
			version:     "a.2.3",
			expectError: true,
		},
		{
			name:        "invalid minor",
			version:     "1.b.3",
			expectError: true,
		},
		{
			name:        "invalid patch",
			version:     "1.2.c",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := ParseSemantic(tt.version)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, version)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, version)
				assert.Equal(t, tt.expected.major, version.Major())
				assert.Equal(t, tt.expected.minor, version.Minor())
				assert.Equal(t, tt.expected.patch, version.Patch())
				assert.Equal(t, tt.expected.prerelease, version.Prerelease())
				assert.Equal(t, tt.expected.metadata, version.Metadata())
				assert.Equal(t, tt.expected.original, version.Original())
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name     string
		version  SemVer
		expected string
	}{
		{
			name: "basic version",
			version: SemVer{
				major: 1,
				minor: 2,
				patch: 3,
			},
			expected: "1.2.3",
		},
		{
			name: "version with prerelease",
			version: SemVer{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha",
			},
			expected: "1.2.3-alpha",
		},
		{
			name: "version with prerelease and metadata",
			version: SemVer{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha.1",
				metadata:   "build.123",
			},
			expected: "1.2.3-alpha.1+build.123",
		},
		{
			name: "version with metadata only",
			version: SemVer{
				major:    1,
				minor:    2,
				patch:    3,
				metadata: "build.123",
			},
			expected: "1.2.3+build.123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.version.String())
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "equal versions",
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: 0,
		},
		{
			name:     "v1 major greater",
			v1:       "2.0.0",
			v2:       "1.9.9",
			expected: 1,
		},
		{
			name:     "v1 major less",
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: -1,
		},
		{
			name:     "v1 minor greater",
			v1:       "1.3.0",
			v2:       "1.2.9",
			expected: 1,
		},
		{
			name:     "v1 minor less",
			v1:       "1.2.0",
			v2:       "1.3.0",
			expected: -1,
		},
		{
			name:     "v1 patch greater",
			v1:       "1.2.4",
			v2:       "1.2.3",
			expected: 1,
		},
		{
			name:     "v1 patch less",
			v1:       "1.2.3",
			v2:       "1.2.4",
			expected: -1,
		},
		{
			name:     "release vs prerelease",
			v1:       "1.2.3",
			v2:       "1.2.3-alpha",
			expected: 1,
		},
		{
			name:     "prerelease vs release",
			v1:       "1.2.3-alpha",
			v2:       "1.2.3",
			expected: -1,
		},
		{
			name:     "prerelease comparison",
			v1:       "1.2.3-beta",
			v2:       "1.2.3-alpha",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := ParseSemantic(tt.v1)
			assert.NoError(t, err)

			v2, err := ParseSemantic(tt.v2)
			assert.NoError(t, err)

			result := v1.Compare(v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_IsPrerelease(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "release version",
			version:  "1.2.3",
			expected: false,
		},
		{
			name:     "prerelease version",
			version:  "1.2.3-alpha",
			expected: true,
		},
		{
			name:     "version with metadata only",
			version:  "1.2.3+build",
			expected: false,
		},
		{
			name:     "prerelease with metadata",
			version:  "1.2.3-alpha+build",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseSemantic(tt.version)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, v.IsPrerelease())
		})
	}
}

func TestVersion_CompatibleWith(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		{
			name:     "same version",
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: true,
		},
		{
			name:     "compatible patch upgrade",
			v1:       "1.2.4",
			v2:       "1.2.3",
			expected: true,
		},
		{
			name:     "compatible minor upgrade",
			v1:       "1.3.0",
			v2:       "1.2.3",
			expected: true,
		},
		{
			name:     "incompatible major upgrade",
			v1:       "2.0.0",
			v2:       "1.2.3",
			expected: false,
		},
		{
			name:     "incompatible downgrade",
			v1:       "1.2.2",
			v2:       "1.2.3",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := ParseSemantic(tt.v1)
			assert.NoError(t, err)

			v2, err := ParseSemantic(tt.v2)
			assert.NoError(t, err)

			result := v1.CompatibleWith(v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_Equal(t *testing.T) {
	v1, err := ParseSemantic("1.2.3")
	assert.NoError(t, err)

	v2, err := ParseSemantic("1.2.3")
	assert.NoError(t, err)

	v3, err := ParseSemantic("1.2.4")
	assert.NoError(t, err)

	assert.True(t, v1.Equal(v2))
	assert.False(t, v1.Equal(v3))
}

func TestVersion_LessThan(t *testing.T) {
	v1, err := ParseSemantic("1.2.3")
	assert.NoError(t, err)

	v2, err := ParseSemantic("1.2.4")
	assert.NoError(t, err)

	assert.True(t, v1.LessThan(v2))
	assert.False(t, v2.LessThan(v1))
}

func TestVersion_GreaterThan(t *testing.T) {
	v1, err := ParseSemantic("1.2.4")
	assert.NoError(t, err)

	v2, err := ParseSemantic("1.2.3")
	assert.NoError(t, err)

	assert.True(t, v1.GreaterThan(v2))
	assert.False(t, v2.GreaterThan(v1))
}
