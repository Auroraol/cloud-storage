export const getConfigFileName = (env) => {
  return `__PRODUCTION__${env.VITE_GLOB_APP_SHORT_NAME || "__APP"}__CONF__`.toUpperCase().replace(/\s/g, "")
}

export function getAppEnvConfig() {
  const ENV_NAME = getConfigFileName(import.meta.env)

  const ENV = import.meta.env.DEV
    ? // 获取全局配置（配置在打包时会独立提取）
      import.meta.env
    : window[ENV_NAME]

  const {
    VITE_GLOB_APP_TITLE,
    VITE_GLOB_API_URL,
    VITE_GLOB_APP_SHORT_NAME,
    VITE_GLOB_API_URL_PREFIX,
    VITE_GLOB_UPLOAD_URL,
    VITE_GLOB_IMG_URL,
    VITE_GLOB_AVATAR_URL
  } = ENV

  if (!/^[a-zA-Z_]*$/.test(VITE_GLOB_APP_SHORT_NAME)) {
    console.log(`VITE_GLOB_APP_SHORT_NAME 变量只能是字母/下划线，请在环境变量中修改并重新运行。`)
  }
  return {
    VITE_GLOB_APP_TITLE,
    VITE_GLOB_API_URL,
    VITE_GLOB_APP_SHORT_NAME,
    VITE_GLOB_API_URL_PREFIX,
    VITE_GLOB_UPLOAD_URL,
    VITE_GLOB_IMG_URL,
    VITE_GLOB_AVATAR_URL
  }
}

/**
 * @description: 开发模式
 */
export const devMode = "development"

/**
 * @description: 生产模式
 */
export const prodMode = "production"

/**
 * @description: 是否为开发模式
 * @returns:
 * @example:
 */
export function isDevMode() {
  return import.meta.env.DEV
}

/**
 * @description: 是否为生产模式
 * @returns:
 * @example:
 */
export function isProdMode() {
  return import.meta.env.PROD
}
