<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { Table } from '@/components/Table'
import { Search } from '@/components/Search'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { ElInput, ElSelect, ElOption, ElTag } from 'element-plus'
import { reactive, ref, unref } from 'vue'
import { getStudentList, createStudent, updateStudent, batchDeleteStudents } from '@/api/student'
import type { Student } from '@/api/student/types'
import { useTable } from '@/hooks/web/useTable'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { pageSize, currentPage } = tableState
    const res = await getStudentList({ pageIndex: unref(currentPage), pageSize: unref(pageSize), ...unref(searchParams) })
    return { list: res.data.list || [], total: res.data.total || 0 }
  },
  fetchDelApi: async () => {
    const res = await batchDeleteStudents(unref(ids).map(Number))
    return !!res
  },
})

const { total, loading, dataList, pageSize, currentPage } = tableState
const { getList, getElTableExpose, delList } = tableMethods

const genderMap: Record<number, string> = { 1: '男', 2: '女' }
const genderOptions = [{ label: '男', value: 1 }, { label: '女', value: 2 }]
const statusMap: Record<string, string> = { unemployed: '未就业', employed: '已就业', graduate_school: '升学', abroad: '出国', startup: '创业' }
const statusOptions = Object.entries(statusMap).map(([k, v]) => ({ label: v, value: k }))
const eduMap: Record<string, string> = { associate: '专科', bachelor: '本科', master: '硕士', doctor: '博士' }
const eduOptions = Object.entries(eduMap).map(([k, v]) => ({ label: v, value: k }))

const crudSchemas = reactive<CrudSchema[]>([
  { field: 'selection', search: { hidden: true }, form: { hidden: true }, table: { type: 'selection' } },
  { field: 'index', label: '#', form: { hidden: true }, search: { hidden: true }, table: { type: 'index', width: 50 } },
  { field: 'keyword', label: '搜索', search: { component: 'Input', componentProps: { placeholder: '姓名/学号' } }, form: { hidden: true }, table: { hidden: true } },
  { field: 'name', label: '姓名', width: 80, search: { hidden: true } },
  { field: 'studentNo', label: '学号', search: { hidden: true }, width: 130 },
  { field: 'gender', label: '性别', width: 55, search: { hidden: true }, table: { slots: { default: (d: any) => <span>{genderMap[d.row.gender] || '-'}</span> } } },
  { field: 'college', label: '学院', minWidth: 130, search: { hidden: true }, table: { slots: { default: (d: any) => <span>{d.row.college?.name || '-'}</span> } } },
  { field: 'grade', label: '年级', width: 70, search: { hidden: true } },
  { field: 'educationLevel', label: '学历', width: 65, search: { hidden: true }, table: { slots: { default: (d: any) => <span>{eduMap[d.row.educationLevel] || d.row.educationLevel || '-'}</span> } } },
  { field: 'expectedCity', label: '期望城市', width: 80, search: { hidden: true } },
  { field: 'expectedSalary', label: '期望薪资', width: 130, search: { hidden: true }, table: { slots: { default: (d: any) => <span>{d.row.expectedSalaryMin && d.row.expectedSalaryMax ? `¥${(d.row.expectedSalaryMin/1000).toFixed(0)}-${(d.row.expectedSalaryMax/1000).toFixed(0)}K` : '-'}</span> } } },
  { field: 'employStatus', label: '就业状态', width: 85, search: { component: 'Select', componentProps: { options: statusOptions, placeholder: '就业状态' } }, table: { slots: { default: (d: any) => <ElTag size="small" type={d.row.employStatus === 'employed' ? 'success' : 'info'}>{statusMap[d.row.employStatus] || d.row.employStatus}</ElTag> } } },
  {
    field: 'action', label: '操作', width: 220, form: { hidden: true }, search: { hidden: true },
    table: { slots: { default: (d: any) => (
          <div class="flex gap-6px">
            <BaseButton size="small" plain onClick={() => viewDetail(d.row)}>详情</BaseButton>
            <BaseButton size="small" type="primary" onClick={() => openEdit(d.row)}>编辑</BaseButton>
            <BaseButton type="danger" size="small" plain onClick={() => delData(d.row)}>删除</BaseButton>
          </div>
        ) } },
  },
])

const { allSchemas } = useCrudSchemas(crudSchemas)
const searchParams = ref({})
const setSearchParams = (p: any) => { currentPage.value = 1; searchParams.value = p; getList() }

const ids = ref<string[]>([])
const delLoading = ref(false)
const delData = async (row?: Student) => {
  const el = await getElTableExpose()
  ids.value = row ? [String(row.id)] : el?.getSelectionRows().map((v: Student) => String(v.id)) || []
  delLoading.value = true
  await delList(unref(ids).length).finally(() => { delLoading.value = false })
}

// ── detail dialog ──
const detailVisible = ref(false)
const currentStudent = ref<Student>()
const viewDetail = (row: Student) => { currentStudent.value = row; detailVisible.value = true }

// ── edit dialog ──
const editVisible = ref(false)
const editLoading = ref(false)
const editForm = reactive<Record<string, any>>({})
const editGender = ref<number>()
const editEdu = ref('')
const editStatus = ref('')
const isCreate = ref(false)

function openCreate() {
  Object.assign(editForm, { name: '', studentNo: '', grade: '', expectedCity: '', expectedIndustry: '', expectedSalaryMin: null, expectedSalaryMax: null, employCompany: '', remark: '' })
  editGender.value = undefined
  editEdu.value = ''
  editStatus.value = 'unemployed'
  isCreate.value = true
  editVisible.value = true
}

function openEdit(row: Student) {
  isCreate.value = false
  Object.assign(editForm, {
    name: row.name, studentNo: row.studentNo, grade: row.grade || '',
    expectedCity: row.expectedCity || '', expectedIndustry: row.expectedIndustry || '',
    expectedSalaryMin: row.expectedSalaryMin, expectedSalaryMax: row.expectedSalaryMax,
    employCompany: row.employCompany || '', remark: row.remark || '',
  })
  editGender.value = row.gender
  editEdu.value = row.educationLevel || ''
  editStatus.value = row.employStatus || 'unemployed'
  currentStudent.value = row
  editVisible.value = true
}

async function handleSave() {
  editLoading.value = true
  try {
    const data = { ...editForm, gender: editGender.value, educationLevel: editEdu.value, employStatus: editStatus.value }
    if (isCreate.value) await createStudent(data)
    else await updateStudent(currentStudent.value!.id, data)
    editVisible.value = false
    getList()
  } finally { editLoading.value = false }
}
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap class="h-full">
      <div class="flex items-center justify-between mb-20px">
        <div><span class="text-18px font-600">学生管理</span><span class="text-13px text-gray-400 ml-12px">共 {{ total }} 人</span></div>
        <div class="flex gap-2"><BaseButton type="primary" @click="openCreate">新增学生</BaseButton><BaseButton :loading="delLoading" type="danger" plain @click="delData()">批量删除</BaseButton></div>
      </div>
      <Search :schema="allSchemas.searchSchema" layout="inline" label-width="70px" @reset="setSearchParams" @search="setSearchParams" />
      <div class="mt-16px">
        <Table v-model:current-page="currentPage" v-model:page-size="pageSize" :columns="allSchemas.tableColumns" :data="dataList" :loading="loading" :show-overflow-tooltip="false" @register="tableRegister" :pagination="{ total }" />
      </div>
    </ContentWrap>

    <!-- Detail dialog -->
    <Dialog v-model="detailVisible" title="学生详情" width="560px">
      <div v-if="currentStudent" class="flex flex-col gap-5">
        <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
          <div class="w-[72px] h-[72px] rounded-full bg-blue-500 flex items-center justify-center shrink-0 shadow-sm">
            <span class="text-white text-2xl font-semibold">{{ currentStudent.name?.charAt(0) }}</span>
          </div>
          <div>
            <h3 class="text-[15px] font-bold text-gray-900">{{ currentStudent.name }}</h3>
            <p class="text-[13px] text-gray-500 mt-0.5">{{ currentStudent.studentNo }}</p>
            <span v-if="currentStudent.employStatus && currentStudent.employStatus !== 'unemployed'" class="inline-block mt-1.5 text-[11px] bg-emerald-50 text-emerald-600 px-2 py-0.5 rounded-full font-medium">{{ statusMap[currentStudent.employStatus] }}</span>
          </div>
        </div>

        <div>
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-blue-500" /><h3 class="text-[13px] font-semibold text-gray-700">基本信息</h3></div>
          <div class="grid grid-cols-2 gap-x-4 gap-y-2.5 text-[13px]">
            <div><span class="text-gray-400">性别</span><p class="text-gray-700 mt-0.5">{{ genderMap[currentStudent.gender] || '-' }}</p></div>
            <div><span class="text-gray-400">学院</span><p class="text-gray-700 mt-0.5">{{ currentStudent.college?.name || '-' }}</p></div>
            <div><span class="text-gray-400">年级</span><p class="text-gray-700 mt-0.5">{{ currentStudent.grade || '-' }}</p></div>
            <div><span class="text-gray-400">学历</span><p class="text-gray-700 mt-0.5">{{ eduMap[currentStudent.educationLevel] || currentStudent.educationLevel || '-' }}</p></div>
          </div>
        </div>

        <div v-if="currentStudent.expectedCity || currentStudent.expectedIndustry || currentStudent.expectedSalaryMin">
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-amber-500" /><h3 class="text-[13px] font-semibold text-gray-700">求职意向</h3></div>
          <div class="grid grid-cols-2 gap-x-4 gap-y-2.5 text-[13px]">
            <div><span class="text-gray-400">期望城市</span><p class="text-gray-700 mt-0.5">{{ currentStudent.expectedCity || '-' }}</p></div>
            <div><span class="text-gray-400">期望行业</span><p class="text-gray-700 mt-0.5">{{ currentStudent.expectedIndustry || '-' }}</p></div>
            <div><span class="text-gray-400">期望薪资</span><p class="text-gray-700 mt-0.5">{{ currentStudent.expectedSalaryMin ? `¥${(currentStudent.expectedSalaryMin/1000).toFixed(0)}K ~ ¥${(currentStudent.expectedSalaryMax!/1000).toFixed(0)}K` : '-' }}</p></div>
          </div>
        </div>

        <div v-if="currentStudent.remark">
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-emerald-500" /><h3 class="text-[13px] font-semibold text-gray-700">备注</h3></div>
          <p class="text-[13px] text-gray-600 leading-relaxed">{{ currentStudent.remark }}</p>
        </div>
      </div>
      <template #footer><BaseButton @click="detailVisible = false">关闭</BaseButton></template>
    </Dialog>

    <!-- Edit dialog -->
    <Dialog v-model="editVisible" :title="isCreate ? '新增学生' : '编辑学生'" width="600px" :close-on-click-modal="false">
      <div class="flex flex-col gap-5">
        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-blue-500" /><h3 class="text-[13px] font-semibold text-gray-700">基本信息</h3></div>
          <div class="space-y-3">
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">姓名 <span class="text-red-400">*</span></label><ElInput v-model="editForm.name" placeholder="学生姓名" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">学号</label><ElInput v-model="editForm.studentNo" placeholder="学号" size="large" /></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">性别</label><ElSelect v-model="editGender" placeholder="选择性别" size="large" :key="'gend-' + editVisible"><ElOption v-for="g in genderOptions" :key="g.value" :label="g.label" :value="g.value" /></ElSelect></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">学历</label><ElSelect v-model="editEdu" placeholder="选择学历" size="large" :key="'edu-' + editVisible"><ElOption v-for="e in eduOptions" :key="e.value" :label="e.label" :value="e.value" /></ElSelect></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">年级</label><ElInput v-model="editForm.grade" placeholder="例如：2024级" size="large" /></div>
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-amber-500" /><h3 class="text-[13px] font-semibold text-gray-700">求职意向</h3></div>
          <div class="space-y-3">
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">期望城市</label><ElInput v-model="editForm.expectedCity" placeholder="例如：深圳" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">期望行业</label><ElInput v-model="editForm.expectedIndustry" placeholder="例如：互联网/IT" size="large" /></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">薪资下限</label><ElInput v-model.number="editForm.expectedSalaryMin" placeholder="例如：8000" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">薪资上限</label><ElInput v-model.number="editForm.expectedSalaryMax" placeholder="例如：15000" size="large" /></div>
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-emerald-500" /><h3 class="text-[13px] font-semibold text-gray-700">就业信息</h3></div>
          <div class="space-y-3">
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">就业状态</label><ElSelect v-model="editStatus" placeholder="选择状态" size="large" :key="'stat-' + editVisible"><ElOption v-for="s in statusOptions" :key="s.value" :label="s.label" :value="s.value" /></ElSelect></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">就业单位</label><ElInput v-model="editForm.employCompany" placeholder="就业单位名称" size="large" /></div>
            </div>
            <div><label class="text-[12px] text-gray-500 mb-1.5 block">备注</label><ElInput v-model="editForm.remark" type="textarea" :rows="2" placeholder="备注信息..." size="large" /></div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2.5">
          <BaseButton size="large" @click="editVisible = false">取消</BaseButton>
          <BaseButton size="large" type="primary" :loading="editLoading" @click="handleSave">{{ isCreate ? '创建学生' : '保存修改' }}</BaseButton>
        </div>
      </template>
    </Dialog>
  </div>
</template>
