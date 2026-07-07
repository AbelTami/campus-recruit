export interface MenuItem {
  id: number
  parentId: number | null
  label: string
  title: string
  path: string
  component: string
  status: number
  children: MenuItem[]
  meta: {
    title: string
    icon: string
    permission: string[]
  }
}

export interface MenuListResponse {
  list: MenuItem[]
}
