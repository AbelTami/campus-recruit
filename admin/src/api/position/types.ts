export interface Position {
  id: number
  enterpriseId: number
  title: string
  industryId: number | null
  city: string
  educationRequirement: string
  experienceRequirement: number | null
  salaryMin: number | null
  salaryMax: number | null
  headcount: number
  description: string
  requirement: string
  welfare: string
  status: number
  enterprise?: { id: number; name: string; shortName: string }
  industry?: { id: number; name: string }
}
