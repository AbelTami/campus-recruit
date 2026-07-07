export interface Student {
  id: number
  studentNo: string
  name: string
  gender: number
  collegeId: number
  majorId: number | null
  grade: string
  educationLevel: string
  graduationYear: number
  expectedCity: string
  expectedSalaryMin: number | null
  expectedSalaryMax: number | null
  expectedIndustry: string
  employStatus: string
  employCompany: string
  employPosition: string
  employSalary: number | null
  remark: string
  createdAt: string
  college?: { id: number; name: string }
  major?: { id: number; name: string }
}

export interface StudentListResponse {
  list: Student[]
  total: number
}
