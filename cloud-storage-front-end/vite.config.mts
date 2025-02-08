import { ConfigEnv, UserConfigExport } from "vite"
import { resolve } from "path"
import vue from "@vitejs/plugin-vue"
import vueJsx from "@vitejs/plugin-vue-jsx"
import { createSvgIconsPlugin } from "vite-plugin-svg-icons"
import svgLoader from "vite-svg-loader"
import UnoCSS from "unocss/vite"
import electron from "vite-electron-plugin"
import { loadViteEnv } from "vite-electron-plugin/plugin"
import { rmSync } from "fs"
import pkg from "./package.json"
import { ElementPlusResolver } from "unplugin-vue-components/resolvers"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import { VueHooksPlusResolver } from "@vue-hooks-plus/resolvers"
import { loadEnv } from "vite" // 导入 loadEnv 函数

/** 清空 dist */
rmSync("dist", { recursive: true, force: true })

/** 配置项文档：https://cn.vitejs.dev/config */
export default ({ mode }: ConfigEnv): UserConfigExport => {
  const viteEnv = loadEnv(mode, process.cwd()) // 获取.env文件里定义的环境变量
  console.log("viteEnv", viteEnv)

  return {
    resolve: {
      alias: {
        /** @ 符号指向 src 目录 */
        "@": resolve(__dirname, "./src")
      }
    },
    server: {
      /** 是否自动打开浏览器 */
      open: false,
      /** 设置 host: true 才可以使用 Network 的形式，以 IP 访问项目 */
      host: pkg.env.host,
      /** 端口号 */
      port: pkg.env.port,
      /** 预热常用文件，提高初始页面加载速度 */
      warmup: {
        clientFiles: ["./src/layouts/**/*.vue"]
      },
      proxy: {
        // 用户中心服务
        [`${viteEnv.VITE_APP_BASE_API}/user_center`]: {
          target: viteEnv.VITE_API_URL_1,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^${viteEnv.VITE_APP_BASE_API}`), "") // 移除 BASE_API 前缀后，转发到用户中心服务
        },
        // 文件上传服务
        [`${viteEnv.VITE_APP_BASE_API}/upload_service`]: {
          target: viteEnv.VITE_API_URL_2,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^${viteEnv.VITE_APP_BASE_API}`), "") // 移除 BASE_API 前缀后，转发到文件上传服务
        },
        // 文件下载服务
        [`${viteEnv.VITE_APP_BASE_API}/download_service`]: {
          target: viteEnv.VITE_API_URL_3,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^${viteEnv.VITE_APP_BASE_API}`), "")
        }
      }
    },
    build: {
      /** 单个 chunk 文件的大小超过 2048KB 时发出警告 */
      chunkSizeWarningLimit: 2048,
      /** 禁用 gzip 压缩大小报告 */
      reportCompressedSize: false,
      rollupOptions: {
        output: {
          /**
           * 分块策略
           * 1. 注意这些包名必须存在，否则打包会报错
           * 2. 如果你不想自定义 chunk 分割策略，可以直接移除这段配置
           */
          manualChunks: {
            vue: ["vue", "vue-router", "pinia"],
            element: ["element-plus", "@element-plus/icons-vue"],
            vxe: ["vxe-table", "vxe-table-plugin-element", "xe-utils"]
          }
        }
      }
    },
    /** 混淆器 */
    esbuild:
      mode === "development"
        ? undefined
        : {
            /** 打包时移除 console.log */
            // pure: ["console.log"],
            /** 打包时移除 debugger */
            drop: ["debugger"],
            /** 打包时移除所有注释 */
            legalComments: "none"
          },
    /** Vite 插件 */
    plugins: [
      vue(),
      vueJsx(),
      // // 引入unplugin-auto-import
      AutoImport({
        //引入vue 自动注册api插件
        imports: ["vue", "vue-router", "vuex"], // 配置需要自动导入的库
        include: [/\.[tj]sx?$/, /\.vue$/, /\.vue\?vue/, /\.md$/],
        dts: "src/auto-import.d.ts", // 自动引入生成api的地址
        resolvers: [VueHooksPlusResolver()] //VueHooks 使用
      }),
      // 引入Element-plus 报错
      // AutoImport({
      //   //plus按需引入
      //   resolvers: [ElementPlusResolver()]
      // }),
      // Components({
      //   //plus按需引入
      //   resolvers: [ElementPlusResolver()],
      //   dts: "src/components.d.ts" //自动引入生成的组件的地址
      // }),
      /** 将 SVG 静态图转化为 Vue 组件 */
      svgLoader({ defaultImport: "url" }),
      /** SVG 插件 */
      createSvgIconsPlugin({
        // Specify the icon folder to be cached
        iconDirs: [resolve(process.cwd(), "./src/assets/icons/svg")],
        // Specify symbolId format
        symbolId: "icon-[dir]-[name]",
        inject: "body-first"
      }),
      /** UnoCSS */
      UnoCSS(),
      electron({
        outDir: "dist",
        include: ["script"],
        transformOptions: { sourcemap: false },
        plugins: [
          {
            name: "remove-comments",
            transform: ({ code }) => {
              let content = code
              // 匹配 块级注释、行级注释、Region注释
              // \s 是匹配所有空白符, 包括换行; \S 非空白符, 不包括换行
              const pattern1 = /\/\*[\s\S]*?\*\/|(\s)+\/\/[\s\S]*?[\n]+/g
              content = content.replaceAll(pattern1, "\n")
              // 匹配 所有空行
              const pattern2 = /^\s*[\r\n]/gm
              content = content.replaceAll(pattern2, "")
              return content
            }
          },
          loadViteEnv()
        ]
      })
    ],
    css: {
      postcss: {
        plugins: [
          {
            postcssPlugin: "internal:charset-removal",
            AtRule: {
              charset: (atRule) => {
                if (atRule.name === "charset") {
                  atRule.remove()
                }
              }
            }
          }
        ]
      }
    },
    clearScreen: false
  }
}
