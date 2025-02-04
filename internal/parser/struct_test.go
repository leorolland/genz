package parser

import (
	"go/ast"
	"reflect"
	"testing"

	"github.com/leorolland/genz/internal/testutils"
	"github.com/leorolland/genz/pkg/models"

	"github.com/google/go-cmp/cmp"
)

func TestParseStructSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode         string
		structName     string
		expectedStruct models.Element
	}{
		"basic struct": {
			goCode: `
			package main

			type A struct {}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
			},
		},
		"struct with one attribute": {
			goCode: `
			package main

			type A struct {
				foo string
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "string", InternalName: "string"},
						Comments: []string{},
					},
				},
			},
		},
		"struct with two attributes": {
			goCode: `
			package main

			type A struct {
				foo string
				bar uint
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "string", InternalName: "string"},
						Comments: []string{},
					},
					{
						Name:     "bar",
						Type:     models.Type{Name: "uint", InternalName: "uint"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with doc": {
			goCode: `
			package main

			type A struct {
				//comment 1
				//comment 2
				foo string
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "string", InternalName: "string"},
						Comments: []string{"comment 1", "comment 2"},
					},
				},
			},
		},
		"attribute with inline comment": {
			goCode: `
			package main

			type A struct {
				foo string // foo
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "string", InternalName: "string"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a slice": {
			goCode: `
			package main

			type B struct {
				foo []string
			}
			`,
			structName: "B",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "[]string", InternalName: "[]string"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo A
			}
			`,
			structName: "B",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "main.A", InternalName: "A"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a slice of named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo []A
			}
			`,
			structName: "B",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "[]main.A", InternalName: "[]A"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a map of named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo map[A]A
			}
			`,
			structName: "B",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "map[main.A]main.A", InternalName: "map[A]A"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a struct containing named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo struct {
					bar []A
					baz string
				}
			}
			`,
			structName: "B",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Attributes: []models.Attribute{
					{
						Name:     "foo",
						Type:     models.Type{Name: "struct{bar []main.A; baz string}", InternalName: "struct{bar []A; baz string}"},
						Comments: []string{},
					},
				},
			},
		},
		"one empty method, value receiver": {
			goCode: `
			package main

			type A struct {}

			func (a A) foo() {}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []models.Type{},
						Returns:           []models.Type{},
						Comments:          []string{},
					},
				},
			},
		},
		"one empty method, comments": {
			goCode: `
			package main

			type A struct {}

			// comment 1
			// comment 2
			func (a A) foo() {}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []models.Type{},
						Returns:           []models.Type{},
						Comments:          []string{"comment 1", "comment 2"},
					},
				},
			},
		},
		"one empty method, pointer receiver": {
			goCode: `
			package main

			type A struct {}

			func (a *A) foo() {}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: true,
						Params:            []models.Type{},
						Returns:           []models.Type{},
						Comments:          []string{},
					},
				},
			},
		},
		"one method with 1 param and 1 return, value receiver": {
			goCode: `
			package main

			type A struct {}

			func (a A) foo(a string) int {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []models.Type{{Name: "string", InternalName: "string"}},
						Returns:           []models.Type{{Name: "int", InternalName: "int"}},
						Comments:          []string{},
					},
				},
			},
		},
		"one method with 1 param and 1 return, named type": {
			goCode: `
			package main

			type T struct {}
			type A struct {}

			func (a A) foo(a T) T {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []models.Type{{Name: "main.T", InternalName: "T"}},
						Returns:           []models.Type{{Name: "main.T", InternalName: "T"}},
						Comments:          []string{},
					},
				},
			},
		},
		"one method with 1 param and 1 return, complex named type": {
			goCode: `
			package main

			type T struct {}
			type A struct {}

			func (a A) foo(a map[T]T) struct{name T} {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []models.Type{{Name: "map[main.T]main.T", InternalName: "map[T]T"}},
						Returns:           []models.Type{{Name: "struct{name main.T}", InternalName: "struct{name T}"}},
						Comments:          []string{},
					},
				},
			},
		},
		"one exported method with 2 params and 2 returns, pointer receiver, comment": {
			goCode: `
			package main

			type A struct {}

			// comment
			func (a *A) Foo(a string, b uint) (int, error) {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type:       models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{},
				Methods: []models.Method{
					{
						Name:              "Foo",
						IsExported:        true,
						IsPointerReceiver: true,
						Params:            []models.Type{{Name: "string", InternalName: "string"}, {Name: "uint", InternalName: "uint"}},
						Returns:           []models.Type{{Name: "int", InternalName: "int"}, {Name: "error", InternalName: "error"}},
						Comments:          []string{"comment"},
					},
				},
			},
		},
		"imported type": {
			goCode: `
			package main

			import "github.com/google/uuid"

			type A struct {
				foo string
				bar uuid.UUID
				baz map[uuid.UUID]uuid.UUID
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{
					{
						Name: "foo",
						Type: models.Type{
							Name:         "string",
							InternalName: "string",
						},
						Comments: []string{},
					},
					{
						Name: "bar",
						Type: models.Type{
							Name:         "uuid.UUID",
							InternalName: "UUID",
						},
						Comments: []string{},
					},
					{
						Name: "baz",
						Type: models.Type{
							Name:         "map[uuid.UUID]uuid.UUID",
							InternalName: "map[UUID]UUID",
						},
						Comments: []string{},
					},
				},
			},
		},
		"struct with tags": {
			goCode: `
			package main

			type A struct {
				foo string ` + "`json:\"foo\"`" + `
				bar string ` + "`json:\"bar\" xml:\"bar\"`" + `
			}
			`,
			structName: "A",
			expectedStruct: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Attributes: []models.Attribute{
					{
						Name: "foo",
						Type: models.Type{
							Name:         "string",
							InternalName: "string",
						},
						Comments: []string{},
						Tags:     map[string]string{"json": "foo"},
					},
					{
						Name: "bar",
						Type: models.Type{
							Name:         "string",
							InternalName: "string",
						},
						Comments: []string{},
						Tags:     map[string]string{"json": "bar", "xml": "bar"},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pkg := testutils.CreatePkgWithCode(t, tc.goCode)

			expr, err := loadAstExpr(pkg, tc.structName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			gotStruct, err := parseStruct(pkg, tc.structName, expr.(*ast.StructType))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(gotStruct, tc.expectedStruct) {
				t.Fatalf("output struct doesn't match expected:\n%s", cmp.Diff(gotStruct, tc.expectedStruct))
			}
		})
	}
}

func Test_parseTags(t *testing.T) {
	testCases := map[string]struct {
		tags    string
		want    map[string]string
		wantErr bool
	}{
		"empty tags": {
			tags:    "",
			want:    map[string]string{},
			wantErr: false,
		},
		"one tag": {
			tags:    "`json:\"name\"`",
			want:    map[string]string{"json": "name"},
			wantErr: false,
		},
		"two tags": {
			tags:    "`json:\"name\" xml:\"name\"`",
			want:    map[string]string{"json": "name", "xml": "name"},
			wantErr: false,
		},
		"tag with options": {
			tags:    "`json:\"name,omitempty\"`",
			want:    map[string]string{"json": "name,omitempty"},
			wantErr: false,
		},
		"tag with options and spaces": {
			tags:    "`json:\"name, omitempty\"`",
			want:    map[string]string{"json": "name, omitempty"},
			wantErr: false,
		},
		"two tags with options": {
			tags:    "`json:\"name,omitempty\" xml:\"name\"`",
			want:    map[string]string{"json": "name,omitempty", "xml": "name"},
			wantErr: false,
		},
		"malformed tag": {
			tags:    "`json:\"name\" xml\"name\"",
			want:    nil,
			wantErr: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := parseTags(tc.tags)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseTags() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseTags() = %v, want %v", got, tc.want)
			}
		})
	}
}
