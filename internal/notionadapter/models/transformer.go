package models

// MdStringObject is just a shortcut for string maps.
type MdStringObject map[string]string

// CustomTransformer is user-provided hook for handling raw blocks.
type CustomTransformer func(data map[string]any) any
