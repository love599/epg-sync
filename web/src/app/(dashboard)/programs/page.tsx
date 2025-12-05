"use client"

import { useState, useEffect } from "react"
import { useToast } from "@/hooks/use-toast"
import api from "@/lib/api"
import ProgramFilters from "./components/ProgramFilters"
import ProgramTable from "./components/ProgramTable"

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

export default function ProgramsPage() {
  const [programs, setPrograms] = useState<Program[]>([])
  const [channels, setChannels] = useState<Channel[]>([])
  const [loading, setLoading] = useState(false)
  const [channelFilter, setChannelFilter] = useState("all")
  const [dateFilter, setDateFilter] = useState<Date | undefined>(new Date())
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPrograms, setTotalPrograms] = useState(0)
  const pageSize = 50
  const { toast } = useToast()

  useEffect(() => {
    loadChannels()
  }, [])

  useEffect(() => {
     searchPrograms()
  }, [channelFilter, dateFilter, currentPage])

  const loadChannels = async () => {
    try {
      const response = await api.get("/admin/channels")
      setChannels(response.data || [])
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "加载失败",
        description: "无法加载频道列表",
      })
    }
  }

  const searchPrograms = async () => {
    try {
      setLoading(true)
      const params: any = {
        page: currentPage,
        page_size: pageSize,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
      }

      if (channelFilter !== "all") {
        params.channel_id = channelFilter
      }

      if (dateFilter) {
        params.date = dateFilter.toISOString().split("T")[0]
      }

      const response = await api.get("/admin/programs/search", { params })
      setPrograms(response.data.items || [])
      setTotalPrograms(response?.data.meta?.total || 0)
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "搜索失败",
        description: error.response?.data?.error || "无法搜索节目",
      })
      setPrograms([])
      setTotalPrograms(0)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = () => {
    setCurrentPage(1)
    searchPrograms()
  }

  const handleReset = () => {
    setChannelFilter("all")
    setDateFilter(new Date())
    setCurrentPage(1)
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">节目管理</h1>
        <p className="text-gray-500 mt-1">搜索和查看EPG节目信息</p>
      </div>

      <ProgramFilters
        channels={channels}
        channelFilter={channelFilter}
        dateFilter={dateFilter}
        loading={loading}
        onChannelChange={setChannelFilter}
        onDateChange={setDateFilter}
        onSearch={handleSearch}
        onReset={handleReset}
      />

      <ProgramTable
        programs={programs}
        channels={channels}
        loading={loading}
        currentPage={currentPage}
        totalPrograms={totalPrograms}
        pageSize={pageSize}
        channelFilter={channelFilter}
        dateFilter={dateFilter}
        onPageChange={setCurrentPage}
      />
    </div>
  )
}
