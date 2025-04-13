/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_PORT: number
  readonly VITE_API_URL_1: string
  readonly VITE_API_URL_2: string
  readonly VITE_API_URL_3: string
  readonly VITE_API_URL_4: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
