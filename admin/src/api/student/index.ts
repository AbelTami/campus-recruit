import request from '@/axios'
import type { Student, StudentListResponse } from './types'

export const getStudentList = (params: Recordable): Promise<IResponse<StudentListResponse>> =>
  request.get({ url: '/admin/students', params })

export const getStudentById = (id: number): Promise<IResponse<Student>> =>
  request.get({ url: `/admin/students/${id}` })

export const createStudent = (data: Partial<Student>): Promise<IResponse<Student>> =>
  request.post({ url: '/admin/students', data })

export const updateStudent = (id: number, data: Partial<Student>): Promise<IResponse<Student>> =>
  request.put({ url: `/admin/students/${id}`, data })

export const deleteStudent = (id: number): Promise<IResponse> =>
  request.delete({ url: `/admin/students/${id}` })

export const batchDeleteStudents = (ids: number[]): Promise<IResponse> =>
  request.post({ url: '/admin/students/delete', data: { ids } })
