"use client"

import { useState, useEffect, useCallback } from "react"
import { Search } from "lucide-react"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useToast } from "@/hooks/use-toast"
import api from "@/lib/api"

interface ChannelMapping {
  id: number
  canonical_id: string
  provider_id: string
  provider_channel_id: string
  confidence: number
  is_verified: boolean
  created_at: string
  updated_at: string
}

export default function MappingsPage() {
  const [mappings, setMappings] = useState<ChannelMapping[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState("")
  const [providerFilter, setProviderFilter] = useState("all")
  const [verifiedFilter, setVerifiedFilter] = useState("all")
  const { toast } = useToast()
  useEffect(() => {
    loadMappings()
  }, [])

  const loadMappings = async () => {
    try {
      setLoading(true)
      const response = await api.get("/admin/channel-mappings")
      setMappings(response.data || [])
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "加载失败",
        description: error.response?.data?.error || "无法加载映射列表",
      })
    } finally {
      setLoading(false)
    }
  }

  const filteredMappings = mappings.filter((mapping) => {
    const matchesSearch =
      mapping.canonical_id.toLowerCase().includes(searchTerm.toLowerCase()) ||
      mapping.provider_channel_id
        .toLowerCase()
        .includes(searchTerm.toLowerCase())

    const matchesProvider =
      providerFilter === "all" || mapping.provider_id === providerFilter

   

    return matchesSearch && matchesProvider
  })

  const providers = Array.from(new Set(mappings.map((m) => m.provider_id)))

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">频道映射</h1>
        <p className="text-gray-500 mt-1">查看频道与各数据源之间的映射关系</p>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <div className="p-6">
          <div className="flex flex-col gap-4">
            <div className="relative">
              <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                placeholder="搜索频道ID、数据源频道ID..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-12 h-12 text-base border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <label className="text-sm font-medium text-gray-700">
                  数据源
                </label>
                <Select
                  value={providerFilter}
                  onValueChange={setProviderFilter}
                >
                  <SelectTrigger className="h-11 border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200">
                    <SelectValue placeholder="选择数据源" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部数据源</SelectItem>
                    {providers.map((provider) => (
                      <SelectItem key={provider} value={provider}>
                        {provider}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {(searchTerm ||
              providerFilter !== "all" ||
              verifiedFilter !== "all") && (
              <div className="flex items-center gap-2 flex-wrap">
                <span className="text-sm text-gray-500">当前筛选:</span>
                {searchTerm && (
                  <span className="inline-flex items-center gap-1 px-3 py-1 bg-blue-50 text-blue-700 rounded-full text-sm">
                    搜索: {searchTerm}
                    <button
                      onClick={() => setSearchTerm("")}
                      className="hover:bg-blue-100 rounded-full p-0.5"
                    >
                      ×
                    </button>
                  </span>
                )}
                {providerFilter !== "all" && (
                  <span className="inline-flex items-center gap-1 px-3 py-1 bg-purple-50 text-purple-700 rounded-full text-sm">
                    数据源: {providerFilter}
                    <button
                      onClick={() => setProviderFilter("all")}
                      className="hover:bg-purple-100 rounded-full p-0.5"
                    >
                      ×
                    </button>
                  </span>
                )}
              
                <button
                  onClick={() => {
                    setSearchTerm("")
                    setProviderFilter("all")
                  }}
                  className="text-sm text-gray-500 hover:text-gray-700 underline"
                >
                  清除全部
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
      {loading ? (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-12">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-500">加载中...</p>
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <div className="border-b border-gray-200 bg-gray-50 px-6 py-3">
            <h2 className="text-sm font-semibold text-gray-700">
              映射列表
              <span className="ml-2 text-gray-500 font-normal">
                (共 {filteredMappings.length} 条记录)
              </span>
            </h2>
          </div>
          <Table>
            <TableHeader>
              <TableRow className="bg-gray-50 hover:bg-gray-50">
                <TableHead className="font-semibold text-gray-700 text-center">
                  频道ID
                </TableHead>
                <TableHead className="font-semibold text-gray-700 text-center">
                  数据源
                </TableHead>
                <TableHead className="font-semibold text-gray-700 text-center">
                  数据源频道ID
                </TableHead>
                <TableHead className="font-semibold text-gray-700 text-center">
                  匹配度
                </TableHead>
                <TableHead className="font-semibold text-gray-700 text-center">
                  更新时间
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredMappings.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-12">
                    <div className="text-gray-400 text-6xl mb-4">🔍</div>
                    <div className="text-gray-500 font-medium">
                      {searchTerm ||
                      providerFilter !== "all" ||
                      verifiedFilter !== "all"
                        ? "没有找到匹配的映射"
                        : "暂无映射数据"}
                    </div>
                    <div className="text-sm text-gray-400 mt-2">
                      {searchTerm ||
                      providerFilter !== "all" ||
                      verifiedFilter !== "all"
                        ? "尝试调整筛选条件"
                        : ""}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                filteredMappings.map((mapping) => (
                  <TableRow
                    key={mapping.id}
                    className="hover:bg-gray-50 transition-colors"
                  >
                    <TableCell className="font-mono text-sm font-medium text-gray-900 text-center">
                      {mapping.canonical_id}
                    </TableCell>
                    <TableCell align="center">
                      <span className="inline-flex items-center px-2.5 py-1 bg-purple-100 text-purple-800 rounded-md text-xs font-medium text-center">
                        {mapping.provider_id || "未知"}
                      </span>
                    </TableCell>
                    <TableCell align="center" className="font-mono text-sm text-gray-700 text-center">
                      {mapping.provider_channel_id}
                    </TableCell>
                    <TableCell >
                      <div className="flex items-center gap-3">
                        <div className="flex-1 bg-gray-200 rounded-full h-2.5 overflow-hidden">
                          <div
                            className="bg-linear-to-r from-green-400 to-green-600 h-2.5 rounded-full transition-all duration-300"
                            style={{
                              width: `${(mapping.confidence / 3.5) * 100}%`,
                            }}
                          />
                        </div>
                        <span className="text-xs font-medium text-gray-700 w-12 text-right">
                          {((mapping.confidence / 3.5) * 100).toFixed(0)}%
                        </span>
                      </div>
                    </TableCell>
                    <TableCell align="center" className="text-sm text-gray-500">
                      {new Date(mapping.updated_at).toLocaleString("zh-CN", {
                        year: "numeric",
                        month: "2-digit",
                        day: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
