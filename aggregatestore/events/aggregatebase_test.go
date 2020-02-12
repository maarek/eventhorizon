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

package events

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

func TestNewAggregateBase(t *testing.T) {
	id := uuid.New().String()
	agg := NewAggregateBase(TestAggregateType, eh.ID(id))
	if agg == nil {
		t.Fatal("there should be an aggregate")
	}
	if agg.AggregateType() != TestAggregateType {
		t.Error("the aggregate type should be correct: ", agg.AggregateType(), TestAggregateType)
	}
	if agg.EntityID() != eh.ID(id) {
		t.Error("the entity ID should be correct: ", agg.EntityID(), eh.ID(id))
	}
	if agg.Version() != 0 {
		t.Error("the version should be 0:", agg.Version())
	}
}

func TestAggregateVersion(t *testing.T) {
	agg := NewAggregateBase(TestAggregateType, eh.ID(uuid.New().String()))
	if agg.Version() != 0 {
		t.Error("the version should be 0:", agg.Version())
	}

	agg.IncrementVersion()
	if agg.Version() != 1 {
		t.Error("the version should be 1:", agg.Version())
	}
}

func TestAggregateEvents(t *testing.T) {
	id := uuid.New().String()
	agg := NewTestAggregate(eh.ID(id))
	timestamp := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	event1 := agg.StoreEvent(TestAggregateEventType, &TestEventData{"event1"}, timestamp)
	if event1.EventType() != TestAggregateEventType {
		t.Error("the event type should be correct:", event1.EventType())
	}
	if !reflect.DeepEqual(event1.Data(), &TestEventData{"event1"}) {
		t.Error("the data should be correct:", event1.Data())
	}
	if !event1.Timestamp().Equal(timestamp) {
		t.Error("the timestamp should not be zero:", event1.Timestamp())
	}
	if event1.Version() != 1 {
		t.Error("the version should be 1:", event1.Version())
	}
	if event1.AggregateType() != TestAggregateType {
		t.Error("the aggregate type should be correct:", event1.AggregateType())
	}
	if event1.AggregateID() != eh.ID(id) {
		t.Error("the aggregate id should be correct:", event1.AggregateID())
	}
	if event1.String() != "TestAggregateEvent@1" {
		t.Error("the string representation should be correct:", event1.String())
	}
	events := agg.Events()
	if len(events) != 1 {
		t.Fatal("there should be one event stored:", len(events))
	}
	if events[0] != event1 {
		t.Error("the stored event should be correct:", events[0])
	}

	event2 := agg.StoreEvent(TestAggregateEventType, &TestEventData{"event1"}, timestamp)
	if event2.Version() != 2 {
		t.Error("the version should be 2:", event2.Version())
	}

	agg.ClearEvents()
	events = agg.Events()
	if len(events) != 0 {
		t.Error("there should be no events stored:", len(events))
	}
	event3 := agg.StoreEvent(TestAggregateEventType, &TestEventData{"event1"}, timestamp)
	if event3.Version() != 1 {
		t.Error("the version should be 1 after clearing uncommitted events (without applying any):", event3.Version())
	}

	agg = NewTestAggregate(eh.ID(uuid.New().String()))
	event1 = agg.StoreEvent(TestAggregateEventType, &TestEventData{"event1"}, timestamp)
	event2 = agg.StoreEvent(TestAggregateEventType, &TestEventData{"event2"}, timestamp)
	events = agg.Events()
	if len(events) != 2 {
		t.Fatal("there should be 2 events stored:", len(events))
	}
	if events[0] != event1 {
		t.Error("the first stored event should be correct:", events[0])
	}
	if events[1] != event2 {
		t.Error("the second stored event should be correct:", events[0])
	}
}

func init() {
	eh.RegisterAggregate(func(id eh.ID) eh.Aggregate {
		return NewTestAggregate(id)
	})

	eh.RegisterEventData(TestAggregateEventType, func() eh.EventData { return &TestEventData{} })
}

const (
	TestAggregateType        eh.AggregateType = "TestAggregate"
	TestAggregateEventType   eh.EventType     = "TestAggregateEvent"
	TestAggregateCommandType eh.CommandType   = "TestAggregateCommand"
)

type TestAggregateCommand struct {
	TestID  eh.ID
	Content string
}

var _ = eh.Command(TestAggregateCommand{})

func (t TestAggregateCommand) AggregateID() eh.ID              { return t.TestID }
func (t TestAggregateCommand) AggregateType() eh.AggregateType { return TestAggregateType }
func (t TestAggregateCommand) CommandType() eh.CommandType     { return TestAggregateCommandType }

type TestEventData struct {
	Content string
}

type TestAggregate struct {
	*AggregateBase
	event eh.Event
}

var _ = Aggregate(&TestAggregate{})

func NewTestAggregate(id eh.ID) *TestAggregate {
	return &TestAggregate{
		AggregateBase: NewAggregateBase(TestAggregateType, eh.ID(id)),
	}
}

func (a *TestAggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	return nil
}

func (a *TestAggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	a.event = event
	return nil
}
