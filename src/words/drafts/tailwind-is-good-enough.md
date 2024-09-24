---
{
  title: 'Tailwind is good enough',
  date: '2024-09-24',
  tags: ['tailwind', 'css'],
}
---

There are a plethora of CSS tools, frameworks, methodologies, runtimes, etc out
there that attempt to wrangle in the beast that is managing CSS in a project in
a variety of ways. 

## using data attribute variant selectors to target different apps/frontends

```js
plugins: [
	plugin(function ({ addVariant }) {
		addVariant("kiosk", ':is([data-rendering-on="kiosk"] &)');
		addVariant("web", ':is([data-rendering-on="web"] &)');
	}),
],
```

## downsides/gripes
- css grid kinda sucks with tailwind
	- to use the power of grid, you just need to use bespoke CSS (i.e. `grid-template-areas` and `grid-area: *`)
	- 
- logical properties support is half-assed
	- why don't you just update `ml-*` / `pt-*` / etc utilities under the hood to use logical properties ??
	- no `-block-` logical property utilities due to (understandable) naming collision with `-bottom` utilities

## best way to use tailwind
- as utility classes / with a custom set of classes / regular css to handle the "composition" and "block" layers of CUBE methodology
- https://cube.fyi
- https://piccalil.li/blog/cube-css/

## references
- [You Don’t Need a CSS Framework](https://www.infoq.com/articles/no-need-css-framework/)
