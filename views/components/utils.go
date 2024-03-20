package components

// Accepts any attribute from a `templ.Attributes` `map[string]any` and verifies
// that it's a non-empty `string`. If it is, it returns the attribute as a
// `string`. If it isn't, it returns the specified `fallback` value instead.
// Should be treated like a nullish coalescing operator in Javascript to set
// sensible default `string` values for an attribute on a component when you
// don't want to require the consumer of the component to always explicitly
// supply each attribute.
//
// Example: always sets a default value of an empty string for the `alt`
// attribute of an `<img>` tag.
//
//	<img src="..." alt={ defaultStrAttr(attrs["alt"], "") } />
func defaultStrAttr(attr any, fallback string) string {
	if attrStr, ok := attr.(string); ok && attrStr != "" {
		return attrStr
	}
	return fallback
}
