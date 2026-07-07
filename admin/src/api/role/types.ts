export interface RoleItem {
  id: number
  roleName: string
  code: string
  status: number
  remark: string
}

export interface RoleListResponse {
  list: RoleItem[]
}
