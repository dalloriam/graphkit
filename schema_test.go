package graphkit

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dalloriam/graphkit/nodes"
)

func TestSchema_ResolveType(t *testing.T) {
	type fields struct {
		types        map[string]*nodes.Block
		RootQuery    string
		RootMutation string
	}

	testType := "HELLO"
	testBlock := &nodes.Block{}

	type args struct {
		typeName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *nodes.Block
		wantErr bool
	}{
		{
			name: "returns correct block when type is present",
			fields: fields{
				types: map[string]*nodes.Block{testType: testBlock},
			},
			args:    args{testType},
			want:    testBlock,
			wantErr: false,
		},
		{
			name: "returns error when type is not present",
			fields: fields{
				types: make(map[string]*nodes.Block),
			},
			args:    args{testType},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				Types:        tt.fields.types,
				RootQuery:    tt.fields.RootQuery,
				RootMutation: tt.fields.RootMutation,
			}
			got, err := s.ResolveType(tt.args.typeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.ResolveType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Schema.ResolveType() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockFileInfo struct {
	name  string
	size  int64
	mode  os.FileMode
	isDir bool
}

func (m mockFileInfo) Name() string      { return m.name }
func (m mockFileInfo) Size() int64       { return m.size }
func (m mockFileInfo) Mode() os.FileMode { return m.mode }
func (mockFileInfo) ModTime() time.Time  { return time.Now() }
func (m mockFileInfo) IsDir() bool       { return m.isDir }
func (mockFileInfo) Sys() interface{}    { return nil }

func TestSchema_processFile(t *testing.T) {
	type testCase struct {
		name string

		path        string
		info        os.FileInfo
		incomingErr error

		wantErr bool
	}

	cases := []testCase{
		{
			name: "returns an error on incoming error",

			path:        "testdata/schema/root.graphql",
			info:        mockFileInfo{},
			incomingErr: errors.New("something terrible happened"),

			wantErr: true,
		},
		{
			name: "returns an error if file doesnt exist",

			path:        "testdata/noexist.graphql",
			info:        mockFileInfo{name: "noexist.graphql"},
			incomingErr: nil,
			wantErr:     true,
		},
		{
			name: "returns no error when file exists",

			path:        "testdata/schema/root.graphql",
			info:        mockFileInfo{name: "root.graphql"},
			incomingErr: nil,
			wantErr:     false,
		},
		{
			name:        "returns no error when file doesnt exist but doesnt end with .graphql",
			path:        "ILOVEYOU.vbs",
			info:        mockFileInfo{name: "ILOVEYOU.vbs"},
			incomingErr: nil,
			wantErr:     false,
		},
	}

	loader := &schemaLoader{}

	for _, c := range cases {
		err := loader.processFile(c.path, c.info, c.incomingErr)
		if (err != nil) != c.wantErr {
			t.Errorf("got err = %v, wanted %v", err, c.wantErr)
		}
	}
}

func Test_loadFromDisk(t *testing.T) {
	t.Run("returns error when processFile fails", func(t *testing.T) {
		_, err := loadFromDisk("non/existent/path/")
		if err == nil {
			t.Errorf("returned no error when processing a non-existent schema")
		}
	})

	t.Run("returns no error when processFile succeeds", func(t *testing.T) {
		_, err := loadFromDisk("testdata/schema/")
		if err != nil {
			t.Errorf("returned an error when processing an existing schema")
		}
	})
}

func TestLoadSchema(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "returns an error when disk loading fails",
			args:    args{"/nonexistent"},
			wantErr: true,
		},
		{
			name:    "returns an error when schema parsing fails",
			args:    args{"testdata/badschema/"},
			wantErr: true,
		},
		{
			name:    "returns no error when schema parsing succeeds",
			args:    args{"testdata/schema/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadSchema(tt.args.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
