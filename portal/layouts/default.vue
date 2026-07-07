<script setup lang="ts">
const { loggedIn } = useUserSession()
const { showExpiredModal } = useAuth()
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col w-full">
    <Navbar />

    <main class="flex-1 pt-16">
      <slot />
    </main>

    <!-- Token expired modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showExpiredModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="showExpiredModal = false" />
          <div class="relative bg-white rounded-2xl shadow-2xl p-8 max-w-sm w-full text-center animate-fadeIn">
            <div class="w-16 h-16 rounded-full bg-amber-50 flex items-center justify-center mx-auto mb-5">
              <Icon name="heroicons:clock" class="w-8 h-8 text-amber-500" />
            </div>
            <h2 class="text-xl font-bold text-gray-900 mb-2">登录已过期</h2>
            <p class="text-gray-500 text-sm mb-6">为了您的账号安全，请重新登录</p>
            <div class="flex gap-3">
              <button @click="showExpiredModal = false" class="flex-1 py-2.5 rounded-lg border border-gray-200 text-gray-600 text-sm font-medium hover:bg-gray-50 transition-colors">稍后再说</button>
              <NuxtLink to="/login" class="flex-1 py-2.5 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 transition-colors" @click="showExpiredModal = false">重新登录</NuxtLink>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <footer class="bg-gray-900 text-gray-300 py-12">
      <div class="max-w-7xl mx-auto px-4 grid grid-cols-1 md:grid-cols-3 gap-8">
        <div>
          <h3 class="text-xl font-bold text-white mb-4">🎓 大学生就业平台</h3>
          <p class="text-gray-400 text-sm">汇聚优质企业，为大学生精准匹配就业机会</p>
        </div>

        <div>
          <h4 class="font-semibold text-white mb-4">快速导航</h4>
          <ul class="space-y-2 text-sm">
            <li><NuxtLink to="/positions" class="hover:text-white transition-colors">职位搜索</NuxtLink></li>
            <li><NuxtLink to="/enterprises" class="hover:text-white transition-colors">企业名录</NuxtLink></li>
            <li v-if="!loggedIn"><NuxtLink to="/login" class="hover:text-white transition-colors">学生登录</NuxtLink></li>
            <li v-if="!loggedIn"><NuxtLink to="/signup" class="hover:text-white transition-colors">立即注册</NuxtLink></li>
          </ul>
        </div>

        <div>
          <h4 class="font-semibold text-white mb-4">关于平台</h4>
          <ul class="space-y-2 text-sm text-gray-400">
            <li>覆盖 24 个学院 · 50 名学生</li>
            <li>20 家合作企业 · 32 个在招职位</li>
            <li class="mt-3 pt-3 border-t border-gray-700">
              &copy; {{ new Date().getFullYear() }} 大学生就业需求分析系统
            </li>
          </ul>
        </div>
      </div>
    </footer>
  </div>
</template>
