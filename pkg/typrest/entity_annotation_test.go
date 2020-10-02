package typrest_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestCreateEntity(t *testing.T) {
	testcases := []struct {
		TestName string
		Annot    *typast.Annot2
		Expected *typrest.Entity
	}{
		{
			TestName: "No tag param",
			Annot: &typast.Annot2{
				Annot: &typast.Annot{
					Decl: &typast.Decl{
						File: typast.File{
							Path: "project/pkg/source.go",
						},
						Type: &typast.StructDecl{
							TypeDecl: typast.TypeDecl{
								Name: "Model",
							},
						},
					},
				},
			},
			Expected: &typrest.Entity{
				Repo:    "ModelRepo",
				Table:   "models",
				Dialect: "",
				DBCtor:  "",
				Target:  "project/pkg",
				Pkg:     "pkg"},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.Expected, typrest.CreateEntity(tt.Annot))
		})
	}
}
