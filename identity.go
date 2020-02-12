// Copyright (c) 2020 - The Event Horizon authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package eventhorizon is a CQRS/ES toolkit.
package eventhorizon

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

type ID string

type IDProvider interface {
	New() string
}

// // ID is the basic identifier interface for all identifier types in eventhorizon.
// type ID interface {
// 	Equal(other interface{}) bool
// 	IsNil() bool
// 	String() string
// }

// // Basic `string` identitifer implementation.
// type BasicID [16]byte

// var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// func randSeq(n int) [16]byte {
// 	b := make([16]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return b
// }

// // Generate a new BasicID.
// func NewID() BasicID {
// 	return randSeq(16)
// }

// // BasicID implements ID.Equal if other is comparable and equal to this.
// func (bid BasicID) Equal(other interface{}) bool {
// 	switch oid := other.(type) {
// 	case BasicID:
// 		return bid.Equal(oid)
// 	default:
// 		return false
// 	}
// }

// // BasicID implements ID.IsNil if the containing id is nil or empty.
// func (bid BasicID) IsNil() bool {
// 	return len(bid) < 0
// }

// // BasicID implements ID.String and returns the inner string.
// func (bid BasicID) String() string {
// 	return *(*string)(unsafe.Pointer(&bid))
// }
