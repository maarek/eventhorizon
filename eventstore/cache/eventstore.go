// Copyright (c) 2020 - The Event Horizon authors
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

package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	eh "github.com/looplab/eventhorizon"
)

// ErrCouldNotRemoveCachedEntry is when a cached entry could be properly removed.
var ErrCouldNotRemoveCachedEntry = errors.New("could not remove cached entry")

// ErrCouldNotSetCacheEntry is when a cached entry can not be saved.
var ErrCouldNotSaveCacheEntry = errors.New("could not save cache entry")

// Cache is an interface for a caching store.
type Cache interface {
	// Set the entry for the specified key.
	Set(key string, entry []byte) error

	// Get the entry for the specified key.
	Get(key string) ([]byte, error)

	// Delete the entry for the specified key.
	Delete(key string) error
}

// EventStore implements the eventhorizon.EventStore interface and holds
// reference to the cache that will be used to store/invalidate events.
type EventStore struct {
	eh.EventStore

	cache Cache
}

// NewEventStore creates a new cached EventStore from an existing store.
func NewEventStore(store eh.EventStore, cache Cache) *EventStore {
	return &EventStore{
		EventStore: store,
		cache:      cache,
	}
}

// Save implements the Save method of the eventhorizon.EventStore interface.
func (s *EventStore) Save(ctx context.Context, events []eh.Event, originalVersion int) error {
	if len(events) == 0 {
		return eh.EventStoreError{
			Err:       eh.ErrNoEventsToAppend,
			Namespace: eh.NamespaceFromContext(ctx),
		}
	}

	// Remove the event from the cache for aggregateID.
	aggregateID := events[0].AggregateID().String()
	_ = s.cache.Delete(aggregateID)

	return s.EventStore.Save(ctx, events, originalVersion)
}

// Load implements the Load method of the eventhorizon.EventStore interface.
func (s *EventStore) Load(ctx context.Context, id uuid.UUID) ([]eh.Event, error) {
	// Load the cached entry and return them if they are found.
	bytes, _ := s.cache.Get(id.String())
	if bytes != nil {
		if results, ok := unmarshalEvents(bytes); ok {
			return results, nil
		}
	}

	// Load the Store entries.
	events, err := s.EventStore.Load(ctx, id)
	if err != nil {
		return events, err
	}

	bytes, err = marshalEvents(events)
	if err != nil {
		return events, eh.EventStoreError{
			Err:       ErrCouldNotSaveCacheEntry,
			Namespace: eh.NamespaceFromContext(ctx),
		}
	}

	err = s.cache.Set(id.String(), bytes)
	if err != nil {
		return events, eh.EventStoreError{
			Err:       ErrCouldNotSaveCacheEntry,
			Namespace: eh.NamespaceFromContext(ctx),
		}
	}

	return events, nil
}

func unmarshalEvents(bytes []byte) ([]eh.Event, bool) {
	var events []eh.Event
	err := json.Unmarshal(bytes, &events)
	if err != nil {
		return nil, false
	}

	return events, true
}

func marshalEvents(events []eh.Event) ([]byte, error) {
	return json.Marshal(events)
}
