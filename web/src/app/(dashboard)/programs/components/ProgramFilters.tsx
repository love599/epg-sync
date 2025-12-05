import { useState } from "react"
import { Search, RefreshCw, ChevronDownIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface Channel {
  id: number
  channel_id: string
  display_name: string
}

interface ProgramFiltersProps {
  channels: Channel[]
  channelFilter: string
  dateFilter: Date | undefined
  loading: boolean
  onChannelChange: (value: string) => void
  onDateChange: (date: Date | undefined) => void
  onSearch: () => void
  onReset: () => void
}

export default function ProgramFilters({
  channels,
  channelFilter,
  dateFilter,
  loading,
  onChannelChange,
  onDateChange,
  onSearch,
  onReset,
}: ProgramFiltersProps) {
  const [dateOpen, setDateOpen] = useState(false)

  const {timeZone, locale} = Intl.DateTimeFormat().resolvedOptions()

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between bg-gray-50 px-6 py-2">
          <h2 className="text-sm font-semibold text-gray-700">节目筛选</h2>
          <Button
            variant="ghost"
            size="sm"
            onClick={onReset}
            className="text-gray-500 hover:text-gray-700"
          >
            <RefreshCw className="h-4 w-4 mr-1" />
            重置
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 px-6">
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">频道</label>
            <Select value={channelFilter} onValueChange={onChannelChange}>
              <SelectTrigger className="h-11 border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200">
                <SelectValue placeholder="选择频道" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部频道</SelectItem>
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

          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">日期</label>
            <Popover open={dateOpen} onOpenChange={setDateOpen}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  id="date"
                  className="w-full h-11 justify-between font-normal border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
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
                    onDateChange(date)
                    setDateOpen(false)
                  }}
                  
                />
              </PopoverContent>
            </Popover>
          </div>
        </div>

        <div className="px-6 py-4">
          <Button
            onClick={onSearch}
            className="w-full h-11 bg-blue-600 hover:bg-blue-700"
            disabled={loading}
          >
            {loading ? (
              <>
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                搜索中...
              </>
            ) : (
              <>
                <Search className="h-4 w-4 mr-2" />
                搜索节目
              </>
            )}
          </Button>
        </div>

        {(channelFilter !== "all" || dateFilter) && (
          <div className="flex items-center gap-2 flex-wrap border-t border-gray-200 px-6 py-4">
            <span className="text-sm text-gray-500">当前筛选:</span>
            {channelFilter !== "all" && (
              <span className="inline-flex items-center gap-1 px-3 py-1 bg-blue-50 text-blue-700 rounded-full text-sm">
                频道:{" "}
                {
                  channels.find((c) => c.channel_id === channelFilter)
                    ?.display_name
                }
                <button
                  onClick={() => onChannelChange("all")}
                  className="hover:bg-blue-100 rounded-full p-0.5"
                >
                  ×
                </button>
              </span>
            )}
            {dateFilter && (
              <span className="inline-flex items-center gap-1 px-3 py-1 bg-purple-50 text-purple-700 rounded-full text-sm">
                日期: {dateFilter.toLocaleDateString("zh-CN")}
                <button
                  onClick={() => onDateChange(undefined)}
                  className="hover:bg-purple-100 rounded-full p-0.5"
                >
                  ×
                </button>
              </span>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
