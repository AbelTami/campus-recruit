export interface ApiResponse<T> { code: number; data: T; message: string }

export interface PaginatedData<T> { list: T[]; total: number }

export interface Position {
  id: number; title: string; enterpriseId: number; city: string
  educationRequirement: string; experienceRequirement: number | null
  salaryMin: number | null; salaryMax: number | null; headcount: number
  description: string; requirement: string; welfare: string; status: number
  enterprise?: Enterprise; industry?: Industry
}

export interface Enterprise {
  id: number; name: string; shortName: string; city: string; scale: string
  nature: string; address: string; website: string; description: string
  contactName: string; contactPhone: string; contactEmail: string; status: number
  logoUrl?: string
  industry?: Industry
}

export interface Industry { id: number; name: string }

export interface Application {
  id: number; status: string; createdAt: string
  position?: Position; enterprise?: Enterprise
}

export interface PortalStats {
  totalStudents: number; totalPositions: number
  totalEnterprises: number; employmentRate: number
}

export interface User { id: number; username: string; nickname: string; avatar: string; roles: string[] }

export interface LoginResponse { access_token: string; refresh_token: string; expires_in: number; userInfo: User }

export interface RegisterRequest { username: string; password: string; name: string; studentNo: string }
