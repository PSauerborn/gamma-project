package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	log "github.com/sirupsen/logrus"
)

var (
	// define custom errors
	ErrInvalidPatch = errors.New("invalid JSON patch operation")
	ErrInvalidJSON  = errors.New("cannot perform json patch: invalid JSON")
)

// function used to perform JSON patch operation on map instance
func PatchJSON(object map[string]interface{},
	operation []map[string]interface{}) (map[string]interface{}, error) {
	// convert operation into JSON format
	patchJson, err := json.Marshal(operation)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert patch operation to JSON: %+v", err))
		return map[string]interface{}{}, ErrInvalidPatch
	}
	// decode JSON patch operation
	patch, err := jsonpatch.DecodePatch(patchJson)
	if err != nil {
		log.Error(fmt.Errorf("unable to parse Json Patch operation: %+v", err))
		return map[string]interface{}{}, ErrInvalidPatch
	}

	// convert metadata to json
	var metaJson []byte
	metaJson, err = json.Marshal(object)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert object to JSON: %+v", err))
		return map[string]interface{}{}, ErrInvalidJSON
	}

	// apply JSON patch operation
	modified, err := patch.Apply(metaJson)
	if err != nil {
		log.Error(fmt.Errorf("unable to apply JSON patch: %+v", err))
		return map[string]interface{}{}, ErrInvalidPatch
	}

	log.Debug(fmt.Sprintf("successfully applied JSON patch to object: %s", modified))
	// convert final JSON string back to interface
	var meta map[string]interface{}
	if err := json.Unmarshal(modified, &meta); err != nil {
		return meta, ErrInvalidJSON
	}
	return meta, nil
}
