/**
This file is licensed under the MIT license provided down below:

MIT License

Copyright (c) 2024-PRESENT Gustavo L. de Mello

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

-------------------------------------------------------------------------------

This file has source code modified from the UnoCSS repository, licensed under the MIT
license. A copy of the original licensed is provided here https://github.com/unocss/unocss/blob/main/LICENSE
and down below:

MIT License

Copyright (c) 2021-PRESENT Anthony Fu <https://github.com/antfu>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */


/**
 * @param {string} str
 * @param {string} [requiredType]
 *
 * @license MIT
 * @author Anthony Fu <https://github.com/antfu>
 */
function bracketWithType(str, requiredType) {
	if (str && str.startsWith('[') && str.endsWith(']')) {
		/** @type {string | undefined } */
		let base
		/** @type {string | undefined } */
		let hintedType

		const bracketTypeRe = /^\[(color|length|size|position|quoted|string):/i
		const match = str.match(bracketTypeRe)
		if (!match) {
			base = str.slice(1, -1)
		}
		else {
			if (!requiredType)
				hintedType = match[1]
			base = str.slice(match[0].length, -1)
		}

		if (!base)
			return

		// test/preset-attributify.test.ts > fixture5
		if (base === '=""')
			return

		if (base.startsWith('--'))
			base = `var(${base})`

		let curly = 0
		for (const i of base) {
			if (i === '[') {
				curly += 1
			}
			else if (i === ']') {
				curly -= 1
				if (curly < 0)
					return
			}
		}
		if (curly)
			return

		switch (hintedType) {
			case 'string': return base
				.replace(/(^|[^\\])_/g, '$1 ')
				.replace(/\\_/g, '_')

			case 'quoted': return base
				.replace(/(^|[^\\])_/g, '$1 ')
				.replace(/\\_/g, '_')
				.replace(/(["\\])/g, '\\$1')
				.replace(/^(.+)$/, '"$1"')
		}

		return base
			.replace(/(url\(.*?\))/g, v => v.replace(/_/g, '\\_'))
			.replace(/(^|[^\\])_/g, '$1 ')
			.replace(/\\_/g, '_')
			.replace(/(?:calc|clamp|max|min)\((.*)/g, (match) => {
				/** @type {string[]} */
				const vars = []
				return match
					.replace(/var\((--.+?)[,)]/g, (match, g1) => {
						vars.push(g1)
						return match.replace(g1, '--un-calc')
					})
					.replace(/(-?\d*\.?\d(?!-\d.+[,)](?![^+\-/*])\D)(?:%|[a-z]+)?|\))([+\-/*])/g, '$1 $2 ')
					.replace(/--un-calc/g, () => /** @type {string} */(vars.shift()))
			})
	}
}

/**
 * @param {string} str
 *
 * @license MIT
 * @author Anthony Fu <https://github.com/antfu>
 */
export function bracket(str) {
	return bracketWithType(str)
}
