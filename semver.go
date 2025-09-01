package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SemVer 表示一个语义版本
type SemVer struct {
	major      uint64
	minor      uint64
	patch      uint64
	prerelease string
	metadata   string
	original   string
}

// 语义版本正则表达式，遵循 semver 2.0.0 规范
var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z\-\.]+))?(?:\+([0-9A-Za-z\-\.]+))?$`)

// ParseSemantic 解析语义版本字符串
func ParseSemantic(version string) (*SemVer, error) {
	if version == "" {
		return nil, fmt.Errorf("version string cannot be empty")
	}

	// 去除前后空格
	version = strings.TrimSpace(version)

	// 匹配语义版本格式
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid semantic version format: %s", version)
	}

	// 解析主版本号
	major, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	// 解析次版本号
	minor, err := strconv.ParseUint(matches[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	// 解析修订版本号
	patch, err := strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	return &SemVer{
		major:      major,
		minor:      minor,
		patch:      patch,
		prerelease: matches[4], // 预发布版本（可选）
		metadata:   matches[5], // 元数据（可选）
		original:   version,
	}, nil
}

// Major 返回主版本号
func (v *SemVer) Major() uint64 {
	return v.major
}

// Minor 返回次版本号
func (v *SemVer) Minor() uint64 {
	return v.minor
}

// Patch 返回修订版本号
func (v *SemVer) Patch() uint64 {
	return v.patch
}

// Prerelease 返回预发布版本标识
func (v *SemVer) Prerelease() string {
	return v.prerelease
}

// Metadata 返回构建元数据
func (v *SemVer) Metadata() string {
	return v.metadata
}

// String 返回规范化的版本字符串
func (v *SemVer) String() string {
	result := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)

	if v.prerelease != "" {
		result += "-" + v.prerelease
	}

	if v.metadata != "" {
		result += "+" + v.metadata
	}

	return result
}

// Original 返回原始版本字符串
func (v *SemVer) Original() string {
	return v.original
}

// Compare 比较两个版本
// 返回值：-1 表示 v < other，0 表示 v == other，1 表示 v > other
func (v *SemVer) Compare(other *SemVer) int {
	// 比较主版本号
	if v.major < other.major {
		return -1
	} else if v.major > other.major {
		return 1
	}

	// 比较次版本号
	if v.minor < other.minor {
		return -1
	} else if v.minor > other.minor {
		return 1
	}

	// 比较修订版本号
	if v.patch < other.patch {
		return -1
	} else if v.patch > other.patch {
		return 1
	}

	// 比较预发布版本
	// 没有预发布版本的优先级高于有预发布版本的
	if v.prerelease == "" && other.prerelease != "" {
		return 1
	} else if v.prerelease != "" && other.prerelease == "" {
		return -1
	} else if v.prerelease != "" && other.prerelease != "" {
		return strings.Compare(v.prerelease, other.prerelease)
	}

	// 版本号完全相同
	return 0
}

// Equal 判断两个版本是否相等
func (v *SemVer) Equal(other *SemVer) bool {
	return v.Compare(other) == 0
}

// LessThan 判断当前版本是否小于另一个版本
func (v *SemVer) LessThan(other *SemVer) bool {
	return v.Compare(other) < 0
}

// GreaterThan 判断当前版本是否大于另一个版本
func (v *SemVer) GreaterThan(other *SemVer) bool {
	return v.Compare(other) > 0
}

// IsPrerelease 判断是否为预发布版本
func (v *SemVer) IsPrerelease() bool {
	return v.prerelease != ""
}

// CompatibleWith 判断是否与另一个版本兼容
// 兼容性规则：主版本号相同，且当前版本不低于目标版本
func (v *SemVer) CompatibleWith(other *SemVer) bool {
	return v.major == other.major && !v.LessThan(other)
}
