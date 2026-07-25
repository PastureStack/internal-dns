package main

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v3"
)

type RecordA struct {
	Ttl    *uint32  `json:"-" yaml:"ttl,omitempty"`
	Answer []string `json:"answer" yaml:"answer"`
}

type RecordCname struct {
	Ttl    *uint32 `json:"-" yaml:"ttl,omitempty"`
	Answer string  `json:"answer" yaml:"answer"`
}

type RecordPtr struct {
	Ttl    *uint32 `json:"-" yaml:"ttl,omitempty"`
	Answer string  `json:"answer" yaml:"answer"`
}

type RecordTxt struct {
	Ttl    *uint32  `json:"-" yaml:"ttl,omitempty"`
	Answer []string `json:"answer" yaml:"answer"`
}

type ClientAnswers struct {
	Search        []string               `json:"search" yaml:"search"`
	Recurse       []string               `json:"recurse" yaml:"recurse"`
	Authoritative []string               `json:"authoritative" yaml:"authoritative"`
	A             map[string]RecordA     `json:"a" yaml:"a"`
	Cname         map[string]RecordCname `json:"cname" yaml:"cname"`
	Ptr           map[string]RecordPtr   `json:"ptr,omitempty" yaml:"ptr,omitempty"`
	Txt           map[string]RecordTxt   `json:"txt,omitempty" yaml:"txt,omitempty"`
}

type clientAnswersWire struct {
	Search              []string               `json:"search" yaml:"search"`
	Recurse             []string               `json:"recurse" yaml:"recurse"`
	Authoritative       *[]string              `json:"authoritative" yaml:"authoritative"`
	LegacyAuthoritative *[]string              `json:"authorative" yaml:"authorative"`
	A                   map[string]RecordA     `json:"a" yaml:"a"`
	Cname               map[string]RecordCname `json:"cname" yaml:"cname"`
	Ptr                 map[string]RecordPtr   `json:"ptr,omitempty" yaml:"ptr,omitempty"`
	Txt                 map[string]RecordTxt   `json:"txt,omitempty" yaml:"txt,omitempty"`
}

func (wire clientAnswersWire) apply(target *ClientAnswers) {
	target.Search = wire.Search
	target.Recurse = wire.Recurse
	target.A = wire.A
	target.Cname = wire.Cname
	target.Ptr = wire.Ptr
	target.Txt = wire.Txt
	target.Authoritative = nil
	if wire.Authoritative != nil {
		target.Authoritative = *wire.Authoritative
	} else if wire.LegacyAuthoritative != nil {
		target.Authoritative = *wire.LegacyAuthoritative
	}
}

func (answers *ClientAnswers) UnmarshalJSON(data []byte) error {
	var wire clientAnswersWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	wire.apply(answers)
	return nil
}

func (answers *ClientAnswers) UnmarshalYAML(value *yaml.Node) error {
	var wire clientAnswersWire
	if err := value.Decode(&wire); err != nil {
		return err
	}
	wire.apply(answers)
	return nil
}

type Answers map[string]ClientAnswers
