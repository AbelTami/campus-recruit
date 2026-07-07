import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { App } from 'vue'
import { Layout } from '@/utils/routerHelper'
import { useI18n } from '@/hooks/web/useI18n'
import { NO_RESET_WHITE_LIST } from '@/constants'

const { t } = useI18n()

export const constantRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/',
    component: Layout,
    redirect: '/employment/dashboard',
    name: 'Root',
    meta: { hidden: true }
  },
  {
    path: '/redirect',
    component: Layout,
    name: 'RedirectWrap',
    children: [
      {
        path: '/redirect/:path(.*)',
        name: 'Redirect',
        component: () => import('@/views/Redirect/Redirect.vue'),
        meta: {}
      }
    ],
    meta: { hidden: true, noTagsView: true }
  },
  {
    path: '/login',
    component: () => import('@/views/Login/Login.vue'),
    name: 'Login',
    meta: { hidden: true, title: t('router.login'), noTagsView: true }
  },
  {
    path: '/personal',
    component: Layout,
    redirect: '/personal/personal-center',
    name: 'Personal',
    meta: { title: t('router.personal'), hidden: true, canTo: true },
    children: [
      {
        path: 'personal-center',
        component: () => import('@/views/Personal/PersonalCenter/PersonalCenter.vue'),
        name: 'PersonalCenter',
        meta: { title: t('router.personalCenter'), hidden: true, canTo: true }
      }
    ]
  },
  {
    path: '/404',
    component: () => import('@/views/Error/404.vue'),
    name: 'NoFind',
    meta: { hidden: true, title: '404', noTagsView: true }
  },
  // ===== 就业管理（直接放 constantRouterMap，避开动态路由时序问题） =====
  {
    path: '/employment',
    component: Layout,
    redirect: '/employment/dashboard',
    name: 'Employment',
    meta: { title: '就业管理', icon: 'ri:bar-chart-box-line', alwaysShow: true },
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/Dashboard/Analysis.vue'),
        name: 'EmploymentDashboard',
        meta: { title: '就业数据大盘', icon: 'ri:dashboard-line', noCache: true, affix: true }
      },
      {
        path: 'salary-analysis',
        component: () => import('@/views/Analysis/SalaryAnalysis.vue'),
        name: 'SalaryAnalysis',
        meta: { title: '薪资分析', icon: 'ri:money-cny-circle-line' }
      },
      {
        path: 'employment-rate',
        component: () => import('@/views/Analysis/EmploymentRate.vue'),
        name: 'EmploymentRate',
        meta: { title: '就业率分析', icon: 'ri:bar-chart-line' }
      },
      {
        path: 'industry-demand',
        component: () => import('@/views/Analysis/IndustryDemand.vue'),
        name: 'IndustryDemand',
        meta: { title: '行业需求分析', icon: 'ri:pie-chart-line' }
      },
      {
        path: 'skill-gap',
        component: () => import('@/views/Analysis/SkillGap.vue'),
        name: 'SkillGap',
        meta: { title: '技能差距分析', icon: 'ri:radar-line' }
      },
      {
        path: 'matching',
        component: () => import('@/views/Analysis/MatchingTool.vue'),
        name: 'MatchingTool',
        meta: { title: '人岗匹配', icon: 'ri:link' }
      },
      {
        path: 'trend-forecast',
        component: () => import('@/views/Analysis/TrendForecast.vue'),
        name: 'TrendForecast',
        meta: { title: '趋势预测', icon: 'ri:line-chart-line' }
      },
      {
        path: 'reports',
        component: () => import('@/views/Analysis/ReportPage.vue'),
        name: 'ReportPage',
        meta: { title: '就业报告', icon: 'ri:article-line' }
      },
      {
        path: 'students',
        component: () => import('@/views/Student/StudentList.vue'),
        name: 'StudentList',
        meta: { title: '学生管理', icon: 'ri:user-line' }
      },
      {
        path: 'enterprises',
        component: () => import('@/views/Enterprise/EnterpriseList.vue'),
        name: 'EnterpriseList',
        meta: { title: '企业管理', icon: 'ri:building-line' }
      },
      {
        path: 'positions',
        component: () => import('@/views/Position/PositionList.vue'),
        name: 'PositionList',
        meta: { title: '职位管理', icon: 'ri:briefcase-line' }
      },
      {
        path: 'applications',
        component: () => import('@/views/Application/ApplicationList.vue'),
        name: 'ApplicationList',
        meta: { title: '投递管理', icon: 'ri:file-list-line' }
      }
    ]
  },
  // ===== 系统管理 =====
  {
    path: '/system',
    component: Layout,
    redirect: '/system/users',
    name: 'System',
    meta: { title: '系统管理', icon: 'ri:settings-line', alwaysShow: true },
    children: [
      {
        path: 'users',
        component: () => import('@/views/Authorization/User/User.vue'),
        name: 'SystemUsers',
        meta: { title: '用户管理', icon: 'ri:user-line' }
      },
      {
        path: 'roles',
        component: () => import('@/views/Authorization/Role/Role.vue'),
        name: 'SystemRoles',
        meta: { title: '角色管理', icon: 'ri:admin-line' }
      },
      {
        path: 'menus',
        component: () => import('@/views/Authorization/Menu/Menu.vue'),
        name: 'SystemMenus',
        meta: { title: '菜单管理', icon: 'ri:menu-line' }
      }
    ]
  }
]

export const asyncRouterMap: AppRouteRecordRaw[] = []

const router = createRouter({
  history: createWebHashHistory(),
  strict: true,
  routes: constantRouterMap as RouteRecordRaw[],
  scrollBehavior: () => ({ left: 0, top: 0 })
})

export const resetRouter = (): void => {
  router.getRoutes().forEach((route) => {
    const { name } = route
    if (name && !NO_RESET_WHITE_LIST.includes(name as string)) {
      router.hasRoute(name) && router.removeRoute(name)
    }
  })
}

export const setupRouter = (app: App<Element>) => {
  app.use(router)
}

export default router
