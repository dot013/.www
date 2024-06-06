import {
	defineConfig,
	presetIcons,
	presetTypography,
	presetUno,
	presetWebFonts,
	transformerDirectives,
	transformerVariantGroup,
} from 'unocss';
import { definePreset } from 'unocss';
import { variantGetParameter } from '@unocss/rule-utils';

import * as utils from './uno-utils.js';

/**
 * Preset based on https://github.com/tailwindlabs/tailwindcss-container-queries
 */
const presetContainers = definePreset(
	/**
	 * @param {{ containers?: Record<string, string | number>}} [options={}]
	 */
	(options = {}) => {
		const defaultContainers = {
			'xs': '20rem' /* 320px */,
			'sm': '24rem' /* 384px */,
			'md': '28rem' /* 448px */,
			'lg': '32rem' /* 512px */,
			'xl': '36rem' /* 576px */,
			'2xl': '42rem' /* 672px */,
			'3xl': '48rem' /* 768px */,
			'4xl': '56rem' /* 896px */,
			'5xl': '64rem' /* 1024px */,
			'6xl': '72rem' /* 1152px */,
			'7xl': '80rem' /* 1280px */,
		};
		options.containers = {
			...defaultContainers,
			...Object.fromEntries(Object.entries(defaultContainers).map(e => [`h-${e[0]}`, e[1]])),
			...Object.fromEntries(Object.entries(defaultContainers).map(e => [`w-${e[0]}`, e[1]])),
			...options.containers,
		};
		/** @type {import('unocss').Preset} */
		const preset = {
			name: 'preset-containers',
			rules: [
				[/^@container(?:\/(\w+))?(?:-(normal|size))?$/, ([, l, v]) => {
					return {
						'container-type': v ?? 'inline-size',
						'container-name': l,
					}
				}],
			],
			variants: [
				{
					name: '@',
					match(matcher, ctx) {
						if (matcher.startsWith('@container')) {
							return matcher
						}
						const variant = variantGetParameter('@', matcher, ctx.generator.config.separators)
						if (variant) {
							const [match, rest, label] = variant;
							const unit = utils.bracket(match);
							console.log(match, rest, label)

							/** @type {string | undefined } */
							let container;
							if (unit?.startsWith("h:")) {
								container = `(min-height: ${unit.replace('h:', '')})`
							} else if (unit?.startsWith("w:") || unit) {
								container = `(min-width: ${unit.replace('w:', '')})`
							} else {
								/** @type {string | number} */
								const size = options.containers?.[match] ?? '';
								container =
									`(${match.startsWith('h-') ? 'min-height' : 'min-width'}: ` +
									`${typeof size === 'number' ? `${size}px` : size})`
							}

							if (!container) {
								return
							}

							let order = (label ? 1000 : 2000) + Object.keys(options.containers ?? {}).indexOf(match);

							return {
								matcher: rest,
								handle: (input, next) => next({
									...input,
									parent: `${input.parent ? `${input.parent} $$` : ''}@container${label ? ` ${label} ` : ''}${container}`,
									parentOrder: order,
								}),
							}
						}
					},
					multiPass: true,
				}
			],
		}
		return preset;
	})

export default defineConfig({
	theme: {
		colors: {
			'black': '#000000',
			'white': '#ffffff',
			'light-gray': '#9c9c9c',
			'gray': '#383838',
			'dark-gray': '#0c0c0c',
			'mauve': '#cba6f7',
			'yellow': '#f9e2af',
		},
	},
	cli: {
		entry: {
			patterns: [
				'./{components,layouts,pages}/*.templ',
				'./static/*.{js,html}',
			],
			outFile: './static/uno.css',
		},
	},
	presets: [
		presetContainers(),
		presetIcons(),
		presetTypography(),
		presetUno(),
		presetWebFonts({
			provider: 'bunny',
			fonts: {
				sans: {
					name: 'Inter',
					provider: 'none'
				},
				cal: {
					name: 'Cal Sans',
					provider: 'none',
				},
				mono: {
					name: 'Fira Code',
					provider: 'none',
				},
			},
		}),
	],
	transformers: [
		transformerDirectives(),
		transformerVariantGroup(),
	],
});
