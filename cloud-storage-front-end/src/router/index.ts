import { type RouteRecordRaw, createRouter } from "vue-router"
import { history, flatMultiLevelRoutes } from "./helper"
import routeSettings from "@/config/route"

const Layouts = () => import("@/layouts/index.vue")

/**
 * 常驻路由
 * 除了 redirect/403/404/login 等隐藏页面，其他页面建议设置 Name 属性
 */
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: "/redirect",
    component: Layouts,
    meta: {
      hidden: true
    },
    children: [
      {
        path: ":path(.*)",
        component: () => import("@/views/redirect/index.vue")
      }
    ]
  },
  {
    path: "/403",
    component: () => import("@/views/error-page/403.vue"),
    meta: {
      hidden: true
    }
  },
  {
    path: "/404",
    component: () => import("@/views/error-page/404.vue"),
    meta: {
      hidden: true
    },
    alias: "/:pathMatch(.*)*"
  },
  {
    path: "/login",
    component: () => import("@/views/login/index.vue"),
    meta: {
      hidden: true
    }
  },
  {
    path: "/",
    component: Layouts,
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        component: () => import("@/views/dashboard/index.vue"),
        name: "Dashboard",
        meta: {
          title: "首页",
          svgIcon: "dashboard",
          affix: true
        }
      }
    ]
  },
  {
    path: "/share",
    component: Layouts,
    redirect: "/share/index",
    children: [
      {
        path: "index",
        component: () => import("@/views/share/index.vue"),
        name: "Share",
        meta: {
          title: "分享",
          svgIcon: "share-alt"
        }
      }
    ]
  },
  {
    path: "/recycle",
    component: Layouts,
    redirect: "/recycle/index",
    children: [
      {
        path: "index",
        component: () => import("@/views/recycle/index.vue"),
        name: "Recycle",
        meta: {
          title: "回收站",
          svgIcon: "recycle (1)"
        }
      }
    ]
  },
  {
    path: "/log",
    component: Layouts,
    redirect: "/log/read",
    name: "Log",
    meta: {
      title: "日志",
      svgIcon: "empire"
    },
    children: [
      {
        path: "read",
        component: () => import("@/views/log/read.vue"),
        name: "Read",
        meta: {
          title: "日志阅读",
          keepAlive: true
        }
      },
      {
        path: "chart",
        component: () => import("@/views/log/chart.vue"),
        name: "LogChart",
        meta: {
          title: "日志图表",
          icon: "chart"
        }
      }
    ]
  },
  {
    path: "/audit",
    component: Layouts,
    meta: { title: "内容审计", svgIcon: "icon_audit" },
    children: [
      {
        path: "log",
        name: "AuditLog",
        component: () => import("@/views/audit/log.vue"),
        meta: { title: "操作日志" }
      },
      {
        path: "statistics",
        name: "AuditStatistics",
        component: () => import("@/views/audit/statistics.vue"),
        meta: { title: "统计分析" }
      }
    ]
  },
  {
    path: "/transfer-records",
    component: Layouts,
    redirect: "/transfer-records/index",
    children: [
      {
        path: "index",
        component: () => import("@/views/transfer/index.vue"),
        meta: {
          title: "传输记录",
          svgIcon: "transfer"
        }
      }
    ]
  },
  {
    path: "/personal-center",
    component: Layouts,
    children: [
      {
        path: "",
        component: () => import("@/views/personal-center/index.vue"),
        meta: {
          title: "个人中心",
          svgIcon: "user"
        }
      }
    ]
  }
]

/**
 * 动态路由
 * 用来放置有权限 (Roles 属性) 的路由
 * 必须带有 Name 属性
 */
export const dynamicRoutes: RouteRecordRaw[] = []

export const router = createRouter({
  history,
  routes: routeSettings.thirdLevelRouteCache ? flatMultiLevelRoutes(constantRoutes) : constantRoutes
})

/** 重置路由 */
export function resetRouter() {
  // 注意：所有动态路由路由必须带有 Name 属性，否则可能会不能完全重置干净
  try {
    router.getRoutes().forEach((route) => {
      const { name, meta } = route
      if (name && meta.roles?.length) {
        router.hasRoute(name) && router.removeRoute(name)
      }
    })
  } catch {
    // 强制刷新浏览器也行，只是交互体验不是很好
    window.location.reload()
  }
}
