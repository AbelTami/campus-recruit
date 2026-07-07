export interface Enterprise {
  id: number
  name: string
  shortName: string
  industryId: number | null
  scale: string
  nature: string
  city: string
  address: string
  website: string
  description: string
  contactName: string
  contactPhone: string
  contactEmail: string
  status: number
  verified: boolean
  logoUrl?: string
  industry?: { id: number; name: string }
}

export interface Industry {
  id: number
  name: string
  code: string
}
