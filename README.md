# 🎓 Campus-Recruit — 校园就业平台

Full-stack campus recruitment platform with student portal, admin panel, and Go API backend.

## Architecture

```
campus-recruit/
├── backend/     # Go API server (Gin + GORM + PostgreSQL)
├── portal/      # Student portal (Nuxt 4 + Vue 3)
├── admin/       # Admin dashboard (Vue 3 + Element Plus)
└── docs/        # Documentation (coming soon)
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go, Gin, GORM, PostgreSQL, Redis, JWT |
| **Portal** | Nuxt 4, Vue 3, TailwindCSS, nuxt-auth-utils |
| **Admin** | Vue 3, Element Plus, ECharts, Pinia, TypeScript |

## Features

- 🔐 JWT authentication with refresh tokens
- 📋 Student job browsing & application tracking
- 🏢 Enterprise CRUD with logo upload
- 📊 Employment analytics dashboard
- 👥 User/role/menu management (RBAC)
- 🌐 i18n support (zh-CN / en / ja)
- 📱 Responsive design

## Quick Start

See individual directories for setup instructions.
