+++
title = 'Event System'
weight = 110
description = "Arche's event system for entity modification."
+++
Arche features an event system that can be used to get notifications about world modifications, namely:

 - Entity creation and removal
 - Component addition, removal and exchange
 - Changes of entity relation targets

Listeners are particularly useful for automating the management
of supplementary data structures that store entities.
They can be used to automatically insert entities into these structures on creation
or component addition, and to remove them on entity or component removal.

The types of interest here are {{< api ecs Listener >}} and {{< api ecs EntityEvent >}}.

## Subscriptions

A listener must subscribe to certain event types.
These are constants of type {{< api "ecs/event" Subscription >}}:

{{< code-func events_test.go TestSubscriptions >}}

Multiple event types can be combined using bit-wise OR (`|`):

{{< code-func events_test.go TestCombineSubscriptions >}}

Some combinations of event types are already defined as {{< api "ecs/event" Subscription >}}.
E.g., to subscribe to all event types, use `event.All`.

Besides subscribing to event types, subscriptions can be restricted to certain component types
that must be affected by the event.
Component type subscriptions are realized using the same {{< api ecs Mask >}} mechanics that filters use. Create masks with {{< api ecs All >}}:

{{< code-func events_test.go TestComponentSubscriptions >}}

A listener with this component subscription would be notified on changes
that are related to the `Position` or to the `Heading` component.

## Builtin listener

An {{< api ecs World >}} can have at most one {{< api ecs Listener >}}.
If required, this listener can be used to dispatch events to sub-listeners.
Such a listener is provided by {{< api listener Dispatch >}}.
In conjunction with {{< api listener Callback >}},
it is already possible to build a sophisticated event system.

In the following example, we compose a {{< api listener Dispatch >}} from two {{< api listener Callback >}}.
The first one listens to all entity creation and entity removal events.
The second one listens to events where a `Position` or a `Heading` is added to an entity.

{{< code-func events_test.go TestListeners >}}

## Custom listeners

Custom listeners can by created by implementing the interface {{< api ecs Listener >}}.
Here is an example of a listener that listens to addition of a `Position` component:

{{< code events_listener_test.go >}}

## EntityEvent

In `Listener.Notify`, as well as in the callback for {{< api listener Callback >}}, we get an {{< api ecs EntityEvent >}} as argument.
It provides all sort of information about the even, like the affected {{< api ecs Entity >}},
event types covered, components added and removed, and more. See the API docs for {{< api ecs EntityEvent >}} for details.
