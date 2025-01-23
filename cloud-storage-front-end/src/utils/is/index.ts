const toString = Object.prototype.toString

/**
 * @description: 判断值是否未某个类型
 */
export const is = (val: any, type: string): boolean => toString.call(val) === `[object ${type}]`

/**
 * @description:  是否为函数
 */
export const isFunction = (val: any): boolean => is(val, "Function") || is(val, "AsyncFunction")

/**
 * @description: 是否为对象
 */
export const isObject = (val: any): boolean => val !== null && is(val, "Object")

/**
 * @description:  是否为字符串
 */
export const isString = (val: any): boolean => is(val, "String")

/**
 * 判断是否 url
 * */
export const isUrl = (url: string): boolean => /^(http|https):\/\//g.test(url)

/**
 * @description:  是否为boolean类型
 */
export const isBoolean = (val: any): boolean => is(val, "Boolean")

/**
 * @description:  是否为数组
 */
export const isArray = (val: any): boolean => val && Array.isArray(val)

/**
 * @description:  是否为数值
 */
export const isNumber = (val: any): boolean => is(val, "Number")
