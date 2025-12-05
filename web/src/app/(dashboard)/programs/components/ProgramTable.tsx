import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

interface Program {
  id: number
  channel_id: string
  title: string
  description: string
  start_time: string
  end_time: string
  provider_id: string
  category: string
  language: string
  rating: string
  episode_number: string
  season_number: string
  created_at: string
}

interface Channel {
  id: number
  channel_id: string
  display_name: string
}

interface ProgramTableProps {
  programs: Program[]
  channels: Channel[]
  loading: boolean
  currentPage: number
  totalPrograms: number
  pageSize: number
  channelFilter: string
  dateFilter: Date | undefined
  onPageChange: (page: number) => void
}

export default function ProgramTable({
  programs,
  channels,
  loading,
  currentPage,
  totalPrograms,
  pageSize,
  channelFilter,
  dateFilter,
  onPageChange,
}: ProgramTableProps) {
  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-12">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-500">加载中...</p>
        </div>
      </div>
    )
  }

  const totalPages = Math.ceil(totalPrograms / pageSize)

  return (
    <>
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="border-b border-gray-200 bg-gray-50 px-6 py-3">
          <h2 className="text-sm font-semibold text-gray-700">
            节目列表
            {totalPrograms > 0 && (
              <span className="ml-2 text-gray-500 font-normal">
                (共 {totalPrograms} 条记录)
              </span>
            )}
          </h2>
        </div>
        <Table>
          <TableHeader>
            <TableRow className="bg-gray-50 hover:bg-gray-50">
              <TableHead className="font-semibold text-gray-700">
                频道
              </TableHead>
              <TableHead className="font-semibold text-gray-700">
                节目标题
              </TableHead>
              <TableHead className="font-semibold text-gray-700">
                开始时间
              </TableHead>
              <TableHead className="font-semibold text-gray-700">
                结束时间
              </TableHead>
              <TableHead className="font-semibold text-gray-700">
                数据源
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {programs.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-12">
                  <div className="text-gray-400 text-6xl mb-4">🔍</div>
                  <div className="text-gray-500 font-medium">
                    {channelFilter !== "all" || dateFilter
                      ? "没有找到匹配的节目"
                      : "请选择搜索条件"}
                  </div>
                  <div className="text-sm text-gray-400 mt-2">
                    {channelFilter !== "all" || dateFilter
                      ? "尝试调整筛选条件"
                      : "选择频道和日期后点击搜索"}
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              programs.map((program) => (
                <TableRow
                  key={program.id}
                  className="hover:bg-gray-50 transition-colors"
                >
                  <TableCell className="font-medium text-gray-900">
                    <div className="flex items-center gap-2">
                      <span>
                        {channels.find(
                          (c) => c.channel_id === program.channel_id
                        )?.display_name || program.channel_id}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="max-w-xs">
                    <div className="font-medium text-gray-900">
                      {program.title}
                    </div>
                  </TableCell>
                  <TableCell className="text-sm text-gray-700">
                    <div className="flex items-center gap-1">
                      <span className="text-green-600">▶</span>
                      {new Date(program.start_time).toLocaleString("zh-CN", {
                        month: "2-digit",
                        day: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </div>
                  </TableCell>
                  <TableCell className="text-sm text-gray-700">
                    <div className="flex items-center gap-1">
                      <span className="text-red-600">■</span>
                      {new Date(program.end_time).toLocaleString("zh-CN", {
                        month: "2-digit",
                        day: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </div>
                  </TableCell>
                  <TableCell className="text-sm">
                    <span className="inline-flex items-center px-2.5 py-1 bg-purple-100 text-purple-800 rounded-md text-xs font-medium">
                      {program.provider_id || "未知"}
                    </span>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {totalPages > 1 && (
        <div className="bg-white border-t border-gray-200 px-6 py-4 rounded-b-lg">
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-500">
              第{" "}
              <span className="font-medium text-gray-900">
                {(currentPage - 1) * pageSize + 1}
              </span>{" "}
              到{" "}
              <span className="font-medium text-gray-900">
                {Math.min(currentPage * pageSize, totalPrograms)}
              </span>{" "}
              条
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => onPageChange(Math.max(1, currentPage - 1))}
                disabled={currentPage === 1}
                className="border-gray-300"
              >
                ← 上一页
              </Button>
              <span className="px-3 py-1 text-sm font-medium text-gray-700 bg-gray-100 rounded">
                {currentPage} / {totalPages}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() =>
                  onPageChange(Math.min(totalPages, currentPage + 1))
                }
                disabled={currentPage === totalPages}
                className="border-gray-300"
              >
                下一页 →
              </Button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}
