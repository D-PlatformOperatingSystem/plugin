// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

//------------------------------------------------------------
// core types for a genesis definition

// GenesisValidator is an initial validator.
type GenesisValidator struct {
	PubKey KeyText `json:"pub_key"`
	Name   string  `json:"name"`
}

// GenesisDoc defines the initial conditions for a tendermint blockchain, in particular its validator set.
type GenesisDoc struct {
	GenesisTime time.Time          `json:"genesis_time"`
	ChainID     string             `json:"chain_id"`
	Validators  []GenesisValidator `json:"validators"`
	AppHash     []byte             `json:"app_hash"`
	AppOptions  interface{}        `json:"app_options,omitempty"`
}

// SaveAs is a utility method for saving GenensisDoc as a JSON file.
func (genDoc *GenesisDoc) SaveAs(file string) error {
	genDocBytes, err := json.Marshal(genDoc)
	if err != nil {
		return err
	}
	return WriteFile(file, genDocBytes, 0644)
}

// ValidatorHash returns the hash of the validator set contained in the GenesisDoc
func (genDoc *GenesisDoc) ValidatorHash() []byte {
	vals := make([]*Validator, len(genDoc.Validators))
	for i, v := range genDoc.Validators {
		if len(v.PubKey.Data) == 0 {
			panic(Fmt("ValidatorHash pubkey of validator[%v] in gendoc is empty", i))
		}
		pubkey, err := PubKeyFromString(v.PubKey.Data)
		if err != nil {
			panic(Fmt("ValidatorHash PubKeyFromBytes failed:%v", err))
		}
		vals[i] = NewValidator(pubkey)
	}
	vset := NewValidatorSet(vals)
	return vset.Hash()
}

// ValidateAndComplete checks that all necessary fields are present
// and fills in defaults for optional fields left empty
func (genDoc *GenesisDoc) ValidateAndComplete() error {

	if genDoc.ChainID == "" {
		return errors.Errorf("Genesis doc must include non-empty chain_id")
	}

	if len(genDoc.Validators) == 0 {
		return errors.Errorf("The genesis file must have at least one validator")
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = time.Now()
	}

	return nil
}

//------------------------------------------------------------
// Make genesis state from file

// GenesisDocFromJSON unmarshalls JSON data into a GenesisDoc.
func GenesisDocFromJSON(jsonBlob []byte) (*GenesisDoc, error) {
	genDoc := GenesisDoc{}
	err := json.Unmarshal(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string) (*GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read GenesisDoc file")
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, errors.Wrap(err, Fmt("Error reading GenesisDoc at %v", genDocFile))
	}
	return genDoc, nil
}
