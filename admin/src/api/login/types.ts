export interface UserLoginType {
  username: string
  password: string
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  roles: string[]
}

export interface LoginResponseData {
  access_token: string
  refresh_token: string
  expires_in: number
  userInfo: UserInfo
}
