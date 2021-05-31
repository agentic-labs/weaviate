package test

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/client/objects"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/test/acceptance/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkipVectorIndex(t *testing.T) {
	// Import a class with vectorizer 'none' and 'skipVectorIndex: true', import
	// objects without vectors.

	t.Run("create schema", func(t *testing.T) {
		createObjectClass(t, &models.Class{
			Class: "TestSkipVectorIndex",
			VectorIndexConfig: map[string]interface{}{
				"skip": true,
			},
			Vectorizer: "none",
			Properties: []*models.Property{
				&models.Property{
					Name:     "name",
					DataType: []string{"text"},
				},
			},
		})
	})

	id := strfmt.UUID("d1d58565-3c9b-4ca6-ac7f-43f739700a1d")

	t.Run("create object", func(t *testing.T) {
		params := objects.NewObjectsCreateParams().WithBody(
			&models.Object{
				ID:         id,
				Class:      "TestSkipVectorIndex",
				Properties: map[string]interface{}{"name": "Jane Doe"},
			})
		_, err := helper.Client(t).Objects.ObjectsCreate(params, nil)
		require.Nil(t, err, "creation should succeed")
	})

	t.Run("get obj by ID", func(t *testing.T) {
		params := objects.NewObjectsGetParams().WithID(id)
		obj, err := helper.Client(t).Objects.ObjectsGet(params, nil)
		require.Nil(t, err, "object can be retrieved by id")

		assert.Equal(t, "Jane Doe", obj.Payload.Properties.(map[string]interface{})["name"].(string))
	})

	t.Run("tear down", func(t *testing.T) {
		deleteObjectClass(t, "TestSkipVectorIndex")
	})
}
