package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid simple path",
			path:    "test.txt",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			path:    "/tmp/test.txt",
			wantErr: false,
		},
		{
			name:    "valid path with subdirectory",
			path:    "subdir/test.txt",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: false,
		},
		{
			name:    "basic path traversal",
			path:    "../test.txt",
			wantErr: true,
			errMsg:  "path traversal detected: ../test.txt",
		},
		{
			name:    "nested path traversal",
			path:    "../../etc/passwd",
			wantErr: true,
			errMsg:  "path traversal detected: ../../etc/passwd",
		},
		{
			name:    "path traversal in middle",
			path:    "some/../other/file.txt",
			wantErr: true,
			errMsg:  "path traversal detected: some/../other/file.txt",
		},
		{
			name:    "path traversal at end",
			path:    "some/path/..",
			wantErr: true,
			errMsg:  "path traversal detected: some/path/..",
		},
		{
			name:    "multiple path traversals",
			path:    "../../../root/.ssh/id_rsa",
			wantErr: true,
			errMsg:  "path traversal detected: ../../../root/.ssh/id_rsa",
		},
		{
			name:    "windows path traversal",
			path:    "..\\..\\windows\\system32",
			wantErr: true,
			errMsg:  "path traversal detected: ..\\..\\windows\\system32",
		},
		{
			name:    "mixed separators",
			path:    "../some\\path/file.txt",
			wantErr: true,
			errMsg:  "path traversal detected: ../some\\path/file.txt",
		},
		{
			name:    "url encoded path traversal",
			path:    "%2e%2e/etc/passwd",
			wantErr: false,
		},
		{
			name:    "unicode path traversal attempt",
			path:    "．．/etc/passwd",
			wantErr: false,
		},
		{
			name:    "legitimate double dots in filename",
			path:    "my..file.txt",
			wantErr: true,
			errMsg:  "path traversal detected: my..file.txt",
		},
		{
			name:    "dots with no slash",
			path:    "..",
			wantErr: true,
			errMsg:  "path traversal detected: ..",
		},
		{
			name:    "current directory reference",
			path:    "./test.txt",
			wantErr: false,
		},
		{
			name:    "hidden file",
			path:    ".hidden",
			wantErr: false,
		},
		{
			name:    "file with dots extension",
			path:    "test.tar.gz",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePath(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
