export interface Application {
  id: number
  studentId: number
  positionId: number
  enterpriseId: number
  status: string
  createdAt: string
  student?: { id: number; name: string; studentNo: string }
  position?: { id: number; title: string }
  enterprise?: { id: number; name: string }
}
