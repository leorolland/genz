package parser

import (
	"go/ast"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leorolland/genz/internal/testutils"
	"github.com/leorolland/genz/pkg/models"
)

func TestParseInterfaceSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode            string
		interfaceName     string
		expectedInterface models.Element
	}{
		"empty interface": {
			goCode: `
			package main

			type A interface {}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type:    models.Type{Name: "main.A", InternalName: "A"},
				Methods: nil,
			},
		},
		"interface with one method": {
			goCode: `
			package main

			type A interface {
				Foo()	
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with two methods": {
			goCode: `
			package main

			type A interface {
				Foo()
				Bar()
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
					{
						Name:              "Bar",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with one method with comments": {
			goCode: `
			package main

			type A interface {
				//Foo does something
				Foo()
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{"Foo does something"},
					},
				},
			},
		},
		"interface with one method with params": {
			goCode: `
			package main
			
			type A interface {
				Foo(a int, b string)
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with one method with returns": {
			goCode: `
			package main

			type A interface {
				Foo() (int, string)
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with a sub interface": {
			goCode: `	
			package main
		
			type A interface {
				Foo() (int, string)
			}
			
			type B interface {
				// A is a sub interface
				A
				Bar()
			}
			`,
			interfaceName: "B",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{" A is a sub interface"},
					},
					{
						Name:              "Bar",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with method with named params": {
			goCode: `
			package main

			type A interface {
				Foo(a int, b string)
			}`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with imported types as params": {
			goCode: `
			package main

			import "github.com/google/uuid"

			type A interface {
				Foo(a uuid.UUID)
			}`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{{Name: "uuid.UUID", InternalName: "UUID"}},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with multi line comments": {
			goCode: `
			package main

			type A interface {
				// Foo does something
				// Foo does something else
				Foo()
			}`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{" Foo does something", " Foo does something else"},
					},
				},
			},
		},
		"interface with unexported method": {
			goCode: `
			package main

			type A interface {
				foo()
			}`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        false,
						Comments:          []string{},
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

			expr, err := loadAstExpr(pkg, tc.interfaceName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			iface, err := parseInterface(pkg, tc.interfaceName, expr.(*ast.InterfaceType))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(iface, tc.expectedInterface) {
				t.Fatalf("output interface doesn't match expected:\n%s", cmp.Diff(iface, tc.expectedInterface))
			}
		})
	}
}
