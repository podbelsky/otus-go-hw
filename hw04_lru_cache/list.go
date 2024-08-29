package hw04lrucache

import "github.com/pkg/errors"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem) error
	MoveToFront(i *ListItem) error
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{Value: v}
	l.linkToFront(i)
	l.len++

	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.back != nil {
		l.back.Next = i
	} else { // empty list
		l.front = i
	}

	l.back = i
	l.len++

	return i
}

func (l *list) Remove(i *ListItem) error {
	if i == nil {
		return errors.New("can't remove nil item")
	}

	before := i.Prev
	after := i.Next

	// see https://en.wikipedia.org/wiki/Doubly_linked_list#Removing_a_node

	if before == nil {
		l.front = after
	} else {
		before.Next = after
	}

	if after == nil {
		l.back = before
	} else {
		after.Prev = before
	}

	l.len--
	i.Prev, i.Next = nil, nil

	return nil
}

func (l *list) MoveToFront(i *ListItem) error {
	if i != nil && l.front == i { // already front
		return nil
	}

	if err := l.Remove(i); err != nil {
		return errors.Wrap(err, "can't move item to front")
	}

	l.linkToFront(i)
	l.len++

	return nil
}

func (l *list) linkToFront(i *ListItem) {
	i.Next = l.front
	i.Prev = nil

	if l.front != nil {
		l.front.Prev = i
	} else { // empty list
		l.back = i
	}

	l.front = i
}

func NewList() List {
	return new(list)
}
