"use client"

import { useState, useEffect } from "react"
import {
  RefreshCw,
  CheckCircle,
  XCircle,
  Clock,
  ChevronDownIcon,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useToast } from "@/hooks/use-toast"
import api from "@/lib/api"

import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

interface Channel {
  id: number
  channel_id: string
  display_name: string
}

interface SyncLog {
  channel_id: string
  channel_name: string
  status: "success" | "failed" | "running"
  message: string
  timestamp: string
}

export default function SyncPage() {
  const [channels, setChannels] = useState<Channel[]>([])
  const [dateFilter, setDateFilter] = useState<Date | undefined>(new Date())
  const [selectedChannel, setSelectedChannel] = useState("CCTV1")
  const [dateOpen, setDateOpen] = useState(false)
  const [syncing, setSyncing] = useState(false)
  const [syncingAll, setSyncingAll] = useState(false)
  const [forceUpdate, setForceUpdate] = useState(false)
  const [syncLogs, setSyncLogs] = useState<SyncLog[]>([])
  const { toast } = useToast()

  const { timeZone, locale } = Intl.DateTimeFormat().resolvedOptions()

  useEffect(() => {
    loadChannels()
    loadSyncLogs()
  }, [])

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

  const loadSyncLogs = () => {
    const logs = localStorage.getItem("syncLogs")
    if (logs) {
      setSyncLogs(JSON.parse(logs))
    }
  }

  const saveSyncLog = (log: SyncLog) => {
    const newLogs = [log, ...syncLogs].slice(0, 50)
    setSyncLogs(newLogs)
    localStorage.setItem("syncLogs", JSON.stringify(newLogs))
  }

  const handleSync = async () => {
    if (syncing) return

    setSyncing(true)
    const startTime = new Date().toISOString()

    try {
      const endpoint = "/admin/epg/sync"
      const params: any = {}

      params.channel_id = selectedChannel

      if (dateFilter) {
        params.start_date = dateFilter.toISOString().split("T")[0]
        params.end_date = dateFilter.toISOString().split("T")[0]
      }

      const response = await api.post(endpoint, null, { params })

      const channelName =
        selectedChannel === "all"
          ? "全部频道"
          : channels.find((c) => c.channel_id === selectedChannel)
              ?.display_name || selectedChannel

      saveSyncLog({
        channel_id: selectedChannel,
        channel_name: channelName,
        status: "success",
        message: response.data?.message || "同步成功",
        timestamp: startTime,
      })

      toast({
        title: "同步成功",
        description: `${channelName} 的EPG数据已更新`,
      })
    } catch (error: any) {
      const channelName =
        selectedChannel === "all"
          ? "全部频道"
          : channels.find((c) => c.channel_id === selectedChannel)
              ?.display_name || selectedChannel

      saveSyncLog({
        channel_id: selectedChannel,
        channel_name: channelName,
        status: "failed",
        message: error.response?.data?.error || "同步失败",
        timestamp: startTime,
      })

      toast({
        variant: "destructive",
        title: "同步失败",
        description: error.response?.data?.error || "无法同步EPG数据",
      })
    } finally {
      setSyncing(false)
    }
  }

  const handleSyncAll = async () => {
    if (syncingAll) return

    setSyncingAll(true)
    const startTime = new Date().toISOString()

    try {
      const response = await api.post("/admin/job/sync", null, {
        params: { force: forceUpdate },
      })

      saveSyncLog({
        channel_id: "all",
        channel_name: "全部频道",
        status: "success",
        message: response.data?.message || "全量同步任务已启动",
        timestamp: startTime,
      })

      toast({
        title: "同步任务已启动",
        description: forceUpdate
          ? "正在强制更新所有频道的EPG数据"
          : "正在同步所有频道的EPG数据",
      })
    } catch (error: any) {
      saveSyncLog({
        channel_id: "all",
        channel_name: "全部频道",
        status: "failed",
        message: error.response?.data?.error || "启动同步任务失败",
        timestamp: startTime,
      })

      toast({
        variant: "destructive",
        title: "启动失败",
        description: error.response?.data?.error || "无法启动全量同步任务",
      })
    } finally {
      setSyncingAll(false)
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "success":
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case "failed":
        return <XCircle className="h-5 w-5 text-red-500" />
      case "running":
        return <RefreshCw className="h-5 w-5 text-blue-500 animate-spin" />
      default:
        return <Clock className="h-5 w-5 text-gray-400" />
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "success":
        return (
          <span className="px-2 py-1 bg-green-100 text-green-800 rounded text-xs">
            成功
          </span>
        )
      case "failed":
        return (
          <span className="px-2 py-1 bg-red-100 text-red-800 rounded text-xs">
            失败
          </span>
        )
      case "running":
        return (
          <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">
            进行中
          </span>
        )
      default:
        return null
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">EPG同步</h1>
        <p className="text-gray-500 mt-1">手动触发EPG数据同步</p>
      </div>



      {/* 全量同步面板 */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-4">全量同步</h2>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 flex-1">
            <input
              type="checkbox"
              id="forceUpdate"
              checked={forceUpdate}
              onChange={(e) => setForceUpdate(e.target.checked)}
              disabled={syncingAll}
              className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <label htmlFor="forceUpdate" className="text-sm font-medium">
              强制更新（忽略已存在的数据）
            </label>
          </div>
          <Button
            onClick={handleSyncAll}
            disabled={syncingAll}
            variant="default"
            size="lg"
            className="min-w-[180px]"
          >
            {syncingAll ? (
              <>
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                同步中...
              </>
            ) : (
              <>
                <RefreshCw className="h-4 w-4 mr-2" />
                同步全部频道
              </>
            )}
          </Button>
        </div>
        <p className="text-sm text-gray-500 mt-3">
          {forceUpdate
            ? "将重新抓取所有频道的EPG数据，包括已存在的数据"
            : "仅同步缺失的EPG数据，跳过已存在的数据"}
        </p>
      </div>

      {/* 单个频道同步面板 */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-4">单频道同步</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="md:col-span-1">
            <label className="block text-sm font-medium mb-2">选择频道</label>
            <Select
              value={selectedChannel}
              onValueChange={setSelectedChannel}
              disabled={syncing}
            >
              <SelectTrigger>
                <SelectValue placeholder="选择要同步的频道" />
              </SelectTrigger>
              <SelectContent>
                {channels.map((channel) => (
                  <SelectItem
                    key={channel.channel_id}
                    value={channel.channel_id}
                  >
                    {channel.display_name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-2">选择日期</label>
            <Popover open={dateOpen} onOpenChange={setDateOpen}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  id="date"
                  className="w-auto  justify-between font-normal border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
                >
                  <span className="flex items-center gap-2">
                    {dateFilter
                      ? dateFilter.toLocaleDateString(locale, {
                          year: "numeric",
                          month: "2-digit",
                          day: "2-digit",
                        })
                      : "选择日期"}
                  </span>
                  <ChevronDownIcon className="h-4 w-4 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="overflow-hidden p-0" align="start">
                <Calendar
                  mode="single"
                  selected={dateFilter}
                  captionLayout="dropdown"
                  timeZone={timeZone}
                  onSelect={(date) => {
                    setDateFilter(date)
                    setDateOpen(false)
                  }}
                />
              </PopoverContent>
            </Popover>
          </div>
          <div className="flex items-end">
            <Button
              onClick={handleSync}
              disabled={syncing}
              className="w-full h-10"
              size="lg"
            >
              {syncing ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  同步中...
                </>
              ) : (
                <>
                  <RefreshCw className="h-4 w-4 mr-2" />
                  开始同步
                </>
              )}
            </Button>
          </div>
        </div>
      </div>
            <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <div className="flex items-start gap-3">
          <div className="text-blue-600 mt-0.5">
            <svg
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div className="flex-1 text-sm text-blue-800">
            <p className="font-medium">同步说明</p>
            <ul className="mt-2 space-y-1 list-disc list-inside">
              <li>频道同步通常需要几秒到几分钟</li>
              <li>同步期间可以查看其他页面,同步会在后台继续进行</li>
              <li>系统会自动按计划同步,手动同步用于立即更新数据</li>
            </ul>
          </div>
        </div>
      </div>
      {/* 统计信息 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-green-100 rounded">
              <CheckCircle className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">成功</div>
              <div className="text-2xl font-bold">
                {syncLogs.filter((log) => log.status === "success").length}
              </div>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-red-100 rounded">
              <XCircle className="h-6 w-6 text-red-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">失败</div>
              <div className="text-2xl font-bold">
                {syncLogs.filter((log) => log.status === "failed").length}
              </div>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gray-100 rounded">
              <Clock className="h-6 w-6 text-gray-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">总计</div>
              <div className="text-2xl font-bold">{syncLogs.length}</div>
            </div>
          </div>
        </div>
      </div>
      {/* 同步历史 */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b">
          <h2 className="text-lg font-semibold">同步历史</h2>
          <p className="text-sm text-gray-500 mt-1">最近的同步记录</p>
        </div>
        <div className="divide-y">
          {syncLogs.length === 0 ? (
            <div className="p-8 text-center text-gray-500">暂无同步记录</div>
          ) : (
            syncLogs.map((log, index) => (
              <div
                key={index}
                className="p-4 hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-start gap-4">
                  <div className="mt-1">{getStatusIcon(log.status)}</div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-1">
                      <span className="font-medium">{log.channel_name}</span>
                      {getStatusBadge(log.status)}
                    </div>
                    <p className="text-sm text-gray-600 mb-1">{log.message}</p>
                    <p className="text-xs text-gray-400">
                      {new Date(log.timestamp).toLocaleString("zh-CN")}
                    </p>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
