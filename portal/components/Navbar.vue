<script setup lang="ts">
const route = useRoute()
const isOpen = ref(false)
const scrolled = ref(false)
const dropdownOpen = ref(false)

const { loggedIn, user, clear } = useUserSession()

const isTransparent = computed(() => route.path === '/' && !scrolled.value)
const isActive = (path: string) => route.path.startsWith(path)

const userName = computed(() => (user.value as any)?.nickname || (user.value as any)?.username || '')
const userInitial = computed(() => userName.value.charAt(0).toUpperCase())
const userRole = computed(() => {
  const roles = (user.value as any)?.roles
  if (!roles?.length) return ''
  return roles.includes('super_admin') ? '管理员' : '学生'
})

function closeDropdown() { dropdownOpen.value = false }
function toggleDropdown() { dropdownOpen.value = !dropdownOpen.value }

async function handleLogout() {
  closeDropdown()
  await clear()
  await navigateTo('/')
}

onMounted(() => {
  scrolled.value = window.scrollY > 20
  window.addEventListener('scroll', () => { scrolled.value = window.scrollY > 20 })
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement
    if (!target.closest('.user-dropdown')) closeDropdown()
  })
})
onUnmounted(() => { window.removeEventListener('scroll', () => {}) })
</script>
<template>
  <header
    class="w-full fixed top-0 left-0 right-0 z-50 transition-all duration-500 border-b"
    :class="isTransparent
      ? 'bg-transparent border-transparent'
      : 'bg-white/80 backdrop-blur-xl border-gray-100 shadow-sm'"
  >
    <nav class="max-w-6xl mx-auto px-4 flex items-center justify-between h-16">
      <div class="flex items-center gap-8">
        <NuxtLink to="/" class="flex items-center gap-2 shrink-0">
          <span class="text-2xl">🎓</span>
          <span class="text-lg font-bold tracking-tight transition-colors" :class="isTransparent ? 'text-white' : 'text-gray-900'">大学生就业平台</span>
        </NuxtLink>

        <div class="hidden md:flex items-center gap-6">
          <NuxtLink v-for="item in [{ path: '/positions', label: '职位搜索' },{ path: '/enterprises', label: '企业名录' }]" :key="item.path" :to="item.path"
            class="group relative text-sm font-medium transition-colors duration-300 py-1"
            :class="isTransparent ? 'text-white/80 hover:text-white' : isActive(item.path) ? 'text-blue-600' : 'text-gray-500 hover:text-gray-800'">
            {{ item.label }}
            <span
              class="absolute -bottom-0.5 left-0 right-0 h-0.5 rounded-full origin-center transition-transform duration-300 ease-out"
              :class="[
                isTransparent ? 'bg-white/60' : isActive(item.path) ? 'bg-blue-600' : 'bg-gray-200',
                isActive(item.path) ? 'scale-x-100' : 'scale-x-0 group-hover:scale-x-100'
              ]"
            />
          </NuxtLink>
        </div>
      </div>

      <div class="hidden md:flex items-center gap-3">
        <template v-if="loggedIn">
          <!-- User dropdown -->
          <div class="user-dropdown relative">
            <button @click="toggleDropdown"
              class="flex items-center gap-2 px-2 py-1.5 -mx-2 rounded-xl transition-all duration-200 cursor-pointer"
              :class="isTransparent ? 'hover:bg-white/10' : 'hover:bg-gray-100/70'">
              <span class="w-7 h-7 rounded-full bg-blue-500 flex items-center justify-center shrink-0">
                <span class="text-white text-[11px] font-semibold">{{ userInitial }}</span>
              </span>
              <span class="text-sm font-medium" :class="isTransparent ? 'text-white' : 'text-gray-700'">{{ userName }}</span>
              <Icon name="heroicons:chevron-down" class="w-3.5 h-3.5 transition-transform duration-200" :class="[isTransparent ? 'text-white/50' : 'text-gray-400', dropdownOpen && 'rotate-180']" />
            </button>

            <Transition name="dropdown">
              <div v-if="dropdownOpen"
                class="absolute right-0 top-full mt-1.5 w-56 bg-white rounded-xl shadow-lg shadow-gray-200/50 border border-gray-100 overflow-hidden z-50">
                <div class="px-4 py-3">
                  <p class="text-sm font-semibold text-gray-900">{{ userName }}</p>
                  <p class="text-xs text-gray-400 mt-0.5">{{ userRole }}</p>
                </div>
                <div class="border-t border-gray-100" />
                <NuxtLink to="/student" @click="closeDropdown"
                  class="flex items-center gap-2.5 px-4 py-2.5 text-[13px] text-gray-600 hover:bg-gray-50 transition-colors">
                  <Icon name="heroicons:user" class="w-4 h-4 text-gray-400" />
                  个人中心
                </NuxtLink>
                <div class="border-t border-gray-100" />
                <button @click="handleLogout"
                  class="flex items-center gap-2.5 w-full px-4 py-2.5 text-[13px] text-gray-500 hover:bg-red-50 hover:text-red-600 transition-colors">
                  <Icon name="heroicons:arrow-right-on-rectangle" class="w-4 h-4 text-gray-400" />
                  退出登录
                </button>
              </div>
            </Transition>
          </div>
        </template>
        <template v-else>
          <NuxtLink to="/login" class="text-sm font-medium transition-colors" :class="isTransparent ? 'text-white/85 hover:text-white' : 'text-gray-600 hover:text-gray-900'">登录</NuxtLink>
          <NuxtLink to="/signup" class="text-sm px-4 py-2 rounded-lg font-medium transition-all"
            :class="isTransparent ? 'bg-white/15 text-white border border-white/30 hover:bg-white/25' : 'bg-blue-600 text-white hover:bg-blue-700 shadow-sm'">注册</NuxtLink>
        </template>
      </div>

      <button @click="isOpen = !isOpen" class="md:hidden">
        <svg class="h-7 w-7" :class="isTransparent ? 'text-white' : 'text-gray-700'" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" viewBox="0 0 24 24">
          <path v-if="!isOpen" d="M4 6h16M4 12h16M4 18h16"/><path v-else d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </nav>

    <Transition name="slide">
      <div v-if="isOpen" class="md:hidden bg-white border-t border-gray-100 shadow-lg">
        <div class="flex flex-col p-4 gap-1">
          <NuxtLink to="/positions" class="px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50" @click="isOpen=false">职位搜索</NuxtLink>
          <NuxtLink to="/enterprises" class="px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50" @click="isOpen=false">企业名录</NuxtLink>
          <hr class="my-2 border-gray-100" />
          <template v-if="loggedIn">
            <div class="flex items-center gap-3 px-2 py-2 mb-1">
              <div class="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center shrink-0">
                <span class="text-white text-sm font-semibold">{{ userInitial }}</span>
              </div>
              <div>
                <p class="text-sm font-semibold text-gray-900">{{ userName }}</p>
                <p class="text-xs text-gray-400">{{ userRole }}</p>
              </div>
            </div>
            <NuxtLink to="/student" class="px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50 flex items-center gap-3" @click="isOpen=false">
              <Icon name="heroicons:user-circle" class="w-4 h-4 text-gray-400" />个人中心
            </NuxtLink>
            <button @click="isOpen=false; handleLogout()" class="px-4 py-3 rounded-lg text-left text-gray-700 hover:bg-red-50 hover:text-red-500 flex items-center gap-3">
              <Icon name="heroicons:arrow-right-on-rectangle" class="w-4 h-4 text-gray-400" />退出登录
            </button>
          </template>
          <template v-else>
            <NuxtLink to="/login" class="px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50" @click="isOpen=false">登录</NuxtLink>
            <NuxtLink to="/signup" class="px-4 py-3 rounded-lg text-center bg-blue-600 text-white font-medium" @click="isOpen=false">注册</NuxtLink>
          </template>
        </div>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
.slide-enter-active { transition: all 0.25s ease-out; }
.slide-leave-active { transition: all 0.2s ease-in; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-8px); }

.dropdown-enter-active { transition: all 0.2s ease-out; }
.dropdown-leave-active { transition: all 0.15s ease-in; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-4px) scale(0.96); }
</style>
